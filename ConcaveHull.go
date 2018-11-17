package ConcaveHull

/**
	Golang implementation of https://github.com/skipperkongen/jts-algorithm-pack/blob/master/src/org/geodelivery/jap/concavehull/SnapHull.java
	which is a Java port of st_concavehull from Postgis 2.0
 */

import (
	"sort"
	"github.com/furstenheim/go-convex-hull-2d"
	"sync"
	"github.com/furstenheim/SimpleRTree"
	"math"
	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/reducers"
)

const DEFAULT_SEGLENGTH = 0.001
type concaver struct {
	rtree * SimpleRTree.SimpleRTree
	seglength float64
	options *Options
	closestPointsMem []closestPoint
	searchItemsMem []searchItem
	flatPointBuffer flatPointsBuffer
	rtreePool *sync.Pool
}
type Options struct {
	Seglength float64
	EstimatedRatioConcaveConvex int // estimated ratio of number of points between concave and convex hull. Will be used to allocate
	ConcaveHullPool *sync.Pool
}

type concaveHullPoolElement struct {
	fpbMem flatPointsBuffer
	closestPointsMem []closestPoint
	searchItemsMem []searchItem
	rtreePool *sync.Pool // This will be passed down to rtree
	pointsCopy FlatPoints
}
func Compute (points FlatPoints) (concaveHull FlatPoints) {
	return ComputeWithOptions(points, nil)
}
func ComputeWithOptions (points FlatPoints, o *Options) (concaveHull FlatPoints) {
	sort.Sort(lexSorter(points))
	return ComputeFromSortedWithOptions(points, o)
}
func ComputeFromSorted (points FlatPoints) (concaveHull FlatPoints) {
	return ComputeFromSortedWithOptions(points, nil)
}

// Compute concave hull from sorted points. Points are expected to be sorted lexicographically by (x,y)
func ComputeFromSortedWithOptions (points FlatPoints, o *Options) (concaveHull FlatPoints) {
	// Create a copy so that convex hull and index can modify the array in different ways
	var pointsCopy FlatPoints
	var rtreeOptions SimpleRTree.Options
	var isConcaveHullPoolElementsSet bool
	var poolEl *concaveHullPoolElement
	var rtreePool *sync.Pool
	if o != nil && o.ConcaveHullPool != nil {
		concaveHullPoolElementsCandidate := o.ConcaveHullPool.Get()
		if concaveHullPoolElementsCandidate != nil {
			isConcaveHullPoolElementsSet = true
			poolEl = concaveHullPoolElementsCandidate.(*concaveHullPoolElement)
			rtreePool = poolEl.rtreePool
		} else {
			rtreePool = &sync.Pool{} // If we are given a pool we want to initiate caching
		}

	}
	if isConcaveHullPoolElementsSet && cap(poolEl.pointsCopy) >= len(points) {
		pointsCopy = poolEl.pointsCopy[0:0]
	} else {
		pointsCopy = make(FlatPoints, 0, len(points))
	}
	pointsCopy = append(pointsCopy, points...)
	rtreeOptions.RTreePool = rtreePool
	rtreeOptions.UnsafeConcurrencyMode = true // we only access from one goroutine at a time
	rtree := SimpleRTree.NewWithOptions(rtreeOptions)
	var wg sync.WaitGroup
	wg.Add(2)
	// Convex hull
	go func () {
		points = go_convex_hull_2d.NewFromSortedArray(points).(FlatPoints)
		wg.Done()
	}()

	func () {
		rtree.LoadSortedArray(SimpleRTree.FlatPoints(pointsCopy))
		wg.Done()
	}()
	wg.Wait()
	var c concaver
	c.seglength = DEFAULT_SEGLENGTH
	if o != nil && o.Seglength != 0 {
		c.seglength = o.Seglength
	}
	c.rtree = rtree
	if isConcaveHullPoolElementsSet {
			c.closestPointsMem = poolEl.closestPointsMem
			c.searchItemsMem = poolEl.searchItemsMem
			c.flatPointBuffer = poolEl.fpbMem
			c.flatPointBuffer.reset()

	} else {
		c.closestPointsMem = make([]closestPoint, 0 , 2)
		c.searchItemsMem = make([]searchItem, 0 , 2)
		estimatedProportionConcave2Convex := 4
		if c.options != nil && c.options.EstimatedRatioConcaveConvex != 0 {
			estimatedProportionConcave2Convex = c.options.EstimatedRatioConcaveConvex
		}
		c.flatPointBuffer = makeFlatPointBuffer(2 * points.Len() * estimatedProportionConcave2Convex)
	}

	result := c.computeFromSorted(points)
	rtree.Destroy() // free resources
	if o != nil && o.ConcaveHullPool != nil {
		o.ConcaveHullPool.Put(
			&concaveHullPoolElement{
				rtreePool: rtreePool,
				searchItemsMem: c.searchItemsMem,
				closestPointsMem: c.closestPointsMem,
				fpbMem: c.flatPointBuffer,
				pointsCopy: pointsCopy,
			},
		)
	}
	return result
}

func (c * concaver) computeFromSorted (convexHull FlatPoints) (concaveHull FlatPoints) {
	// degerated case
	if (convexHull.Len() < 3) {
		return convexHull
	}

	x0, y0 := convexHull.Take(0)
	concaveHullBuffer := c.flatPointBuffer
	concaveHullBuffer.addFloat(x0)
	concaveHullBuffer.addFloat(y0)
	for i := 0; i<convexHull.Len(); i++ {
		x1, y1 := convexHull.Take(i)
		var x2, y2 float64
		if i == convexHull.Len() -1 {
			x2, y2 = convexHull.Take(0)
		} else {
			x2, y2 = convexHull.Take(i + 1)
		}
		sideSplit := c.segmentize(x1, y1, x2, y2)
		for _, p := range(sideSplit) {
			concaveHullBuffer.addFloat(p.x)
			concaveHullBuffer.addFloat(p.y)
		}
	}
	concaveHull = concaveHullBuffer.toFloatArray()
	path := reducers.DouglasPeucker(geo.NewPathFromFlatXYData(concaveHull), c.seglength)
	// reused allocated array
	concaveHull = concaveHull[0:0]
	reducedPoints := path.Points()

	for _, p := range(reducedPoints) {
		concaveHull = append(concaveHull, p.Lng(), p.Lat())
	}
	return concaveHull
}

// Split side in small edges, for each edge find closest point. Remove duplicates
func (c * concaver) segmentize (x1, y1, x2, y2 float64) (points []closestPoint) {
	dist := math.Sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2))
	nSegments := math.Ceil(dist / c.seglength)
	factor := 1 / nSegments
	vX := factor * (x2 - x1)
	vY := factor * (y2 - y1)

	closestPoints := c.closestPointsMem[0: 0]
	closestPoints = append(closestPoints, closestPoint{index: 0, x: x1, y: y1})
	closestPoints = append(closestPoints, closestPoint{index: int(nSegments), x: x2, y: y2})

	if (nSegments < 2) {
		return closestPoints[1:]
	}

	stack := c.searchItemsMem[0: 0]
	stack = append(stack, searchItem{left: 0, right: int(nSegments), lastLeftIndex: 0, lastRightIndex: 1})
	for len(stack) > 0 {
		var item searchItem
		item, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if item.right - item.left <= 1 {
			continue
		}
		index := (item.left + item.right) / 2
		fIndex := float64(index)
		currentX := x1 + vX * fIndex
		currentY := y1 + vY * fIndex
		lx := closestPoints[item.lastLeftIndex].x
		ly := closestPoints[item.lastLeftIndex].y
		rx := closestPoints[item.lastRightIndex].x
		ry := closestPoints[item.lastRightIndex].y

		d1 := (currentX - lx) * (currentX - lx) + (currentY - ly) * (currentY - ly)
		d2 := (currentX - rx) * (currentX - rx) + (currentY - ry) * (currentY - ry)
		x, y, _, found := c.rtree.FindNearestPointWithin(currentX, currentY, math.Min(d1, d2))
		if !found {
			continue
		}
		isNewLeft := x != lx || y != ly
		isNewRight := x != rx || y != ry

		// we don't know the point
		if isNewLeft && isNewRight {
			newResultIndex := len(closestPoints)
			closestPoints = append(closestPoints, closestPoint{index: index, x: x, y: y})
			stack = append(stack, searchItem{left: item.left, right: index, lastLeftIndex: item.lastLeftIndex, lastRightIndex: newResultIndex})
			// alloc
			stack = append(stack, searchItem{left: index, right: item.right, lastLeftIndex: newResultIndex, lastRightIndex: item.lastRightIndex})
		} else if (isNewLeft) {
			stack = append(stack, searchItem{left: item.left, right: index, lastLeftIndex: item.lastLeftIndex, lastRightIndex: item.lastRightIndex})
		} else {
			// don't add point to closest points, but we need to keep looking on the right side
			stack = append(stack, searchItem{left: index, right: item.right, lastLeftIndex: item.lastLeftIndex, lastRightIndex: item.lastRightIndex})
		}
	}
	closestPointSorter(closestPoints).cpSort()
	c.searchItemsMem = stack
	c.closestPointsMem = closestPoints
	return closestPoints[1:]
}

type closestPoint struct {
	index int
	x, y float64
}

type searchItem struct {
	left, right, lastLeftIndex, lastRightIndex int

}





type FlatPoints []float64

func (fp FlatPoints) Len () int {
	return len(fp) / 2
}

func (fp FlatPoints) Slice (i, j int) (go_convex_hull_2d.Interface) {
	return fp[2 * i: 2 * j]
}

func (fp FlatPoints) Swap (i, j int) {
	fp[2 * i], fp[2 * i + 1], fp[2 * j], fp[2 * j + 1] = fp[2 * j], fp[2 * j + 1], fp[2 * i], fp[2 * i + 1]
}

func (fp FlatPoints) Take(i int) (x1, y1 float64) {
	return fp[2 * i], fp[2 * i +1]
}

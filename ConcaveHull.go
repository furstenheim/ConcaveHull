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
)


type concaver struct {
	rtree * SimpleRTree.SimpleRTree
	seglength float64
}
func Compute (points FlatPoints) (concaveHull FlatPoints) {
	sort.Sort(lexSorter(points))
	return ComputeFromSorted(points)
}

// Compute concave hull from sorted points. Points are expected to be sorted lexicographically by (x,y)
func ComputeFromSorted (points FlatPoints) (concaveHull FlatPoints) {
	// Create a copy so that convex hull and index can modify the array in different ways
	pointsCopy := make(FlatPoints, 0, len(points))
	pointsCopy = append(pointsCopy, points...)
	rtree := SimpleRTree.New()
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
	var c concaver
	c.seglength = 0.0001 // TODO get from options
	c.rtree = rtree
	return c.computeFromSorted(points)
}

func (c * concaver) computeFromSorted (convexHull FlatPoints) (concaveHull FlatPoints) {
	// TODO treat degenerated cases of convexHull
	concaveHull = make([]float64, 0, 2 * convexHull.Len())
	for i := 0; i<convexHull.Len(); i++ {
		x1, y1 := convexHull.Take(i)
		var x2, y2 float64
		if i == convexHull.Len() -1 {
			x2, y2 = convexHull.Take(0)
		} else {
			x2, y2 = convexHull.Take( i + 1)
		}
		sideSplit := c.segmentize(x1, y1, x2, y2)
		concaveHull = append(concaveHull, sideSplit...)
	}
	return concaveHull
}

// Split side in small edges, for each edge find closest point. Remove duplicates
func (c * concaver) segmentize (x1, y1, x2, y2 float64) (points []float64) {
	dist := math.Sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2))
	nSegments := math.Ceil(dist / c.seglength)
	factor := 1 / nSegments
	flatPoints := make([]float64, 0, int(2 * nSegments))
	vX := factor * (x2 - x1)
	vY := factor * (y2 - y1)
	flatPoints = append(flatPoints, x1, y1)

	currentX := x1
	currentY := y1

	latestX := x1
	latestY := y1

	for i := 0; i < int(nSegments); i++ {
		x, y, _, _ := c.rtree.FindNearestPoint(currentX, currentY)
		if x != latestX && y != latestY {
			flatPoints = append(flatPoints, x, y)
			latestX = x
			latestY = y
		}
		currentX += vX
		currentY += vY
	}
	return flatPoints
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
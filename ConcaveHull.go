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
)

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
	return points
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
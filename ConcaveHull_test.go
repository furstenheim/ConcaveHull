package ConcaveHull

import (
	"testing"
	"fmt"
)

func TestCompute_convexHullInAntiClockwiseOrder(t *testing.T) {
	points := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	points2 := append([]float64{}, points...)
	points2 = append(points2)
	result := ComputeFromSorted(points)
	fmt.Println(result)
	fmt.Println(points2)
	compareConcaveHulls(t, result, points2)
}


func compareConcaveHulls(t *testing.T, actualC, expectedC FlatPoints) {
	if actualC.Len() != expectedC.Len() {
		t.Errorf("Concave hull didn't correct length, got %d, want: %d", len(actualC), len(expectedC))
		for i := 0; i < actualC.Len(); i++ {
			t.Log(actualC.Take(i))
		}
		return
	}
	for i := 0; i < actualC.Len(); i++ {
		x1, y1 := actualC.Take(i)
		x2, y2 := expectedC.Take(i)
		if x1 != x2 || y1 != y2 {
			fmt.Println(actualC, expectedC)
			t.Errorf("%d th point of the convex hull was not correct, got: %+v want: %+v", i, x1, y1, x2, y2)
		}
	}
}
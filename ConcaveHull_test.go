package ConcaveHull

import (
	"testing"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/paulmach/go.geo"
)

func TestCompute_convexHullInAntiClockwiseOrder(t *testing.T) {
	points := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	points2 := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0., 0.}
	result := ComputeFromSorted(points)
	fmt.Println(result)
	fmt.Println(points2)
	compareConcaveHulls(t, result, points2)
}


func TestCompute_convexHullShuffled(t *testing.T) {
	points := FlatPoints{1.0, 0.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	points2 := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0., 0.}
	result := Compute(points)
	fmt.Println("result", result)
	fmt.Println(points2)
	compareConcaveHulls(t, result, points2)
}

func TestComputeFromSorted_simpleConvexHull(t *testing.T) {
	points := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	points2 := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0., 0.}
	result := ComputeFromSorted(points)
	fmt.Println(result)
	fmt.Println(points2)
	compareConcaveHulls(t, result, points2)
}


func TestCompute_simpleConcaveHull(t *testing.T) {
	points := FlatPoints{1./3., 0.5, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0}
	points2 := FlatPoints{0.0, 0.0, 1.0, 0., 1., 1., 0., 1., 1./3., 0.5, 0.0, 0.0}
	result := Compute(points)
	fmt.Println(result)
	fmt.Println(points2)
	compareConcaveHulls(t, result, points2)
}

func BenchmarkCompute_wkb(b * testing.B) {
	dat, err := ioutil.ReadFile("./wkb")
	if (err != nil) {
		log.Fatal(err)
	}
	path := geo.NewPathFromWKB(dat)
	points := make([]float64, 2 * path.Length())
	for _, p := range(path.Points()) {
		points = append(points, p.Lng(), p.Lat())
	}
	fmt.Println("Length of polygon", path.Length())

	b.Run("Wkb", func (b * testing.B) {
		result := Compute(points)
		fmt.Println(result.Len())
	})


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
			t.Errorf("%d th point of the convex hull was not correct, got: %+v want: %+v", i, x1, y1)
		}
	}
}
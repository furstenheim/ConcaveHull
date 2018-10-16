package ConcaveHull

import (
	"testing"
	"fmt"
	"log"
	"math"
	"math/rand"
	"github.com/furstenheim/SimpleRTree"
	"bufio"
	"os"
	"strings"
	"strconv"
	"path/filepath"
	"github.com/stretchr/testify/assert"
)

func TestCompute_concaveHullInAntiClockwiseOrder(t *testing.T) {
	points := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	points2 := FlatPoints{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0., 0.}
	result := ComputeFromSorted(points)
	fmt.Println(result)
	fmt.Println(points2)
	compareConcaveHulls(t, result, points2)
}


func TestCompute_concaveHullShuffled(t *testing.T) {
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

func TestConcaveHull_segmentize (t *testing.T) {
	const size = 200
	points := make([]float64, size * 2)
	for i := 0; i < 2 * size; i++ {
		points[i] = rand.Float64()
	}
	fp := FlatPoints(points)
	r := SimpleRTree.New().Load(SimpleRTree.FlatPoints(fp))
	c := new(concaver)
	c.rtree = r


	c.seglength = DEFAULT_SEGLENGTH

	for i := 0; i < 1; i++ {
		index1 := rand.Intn(size)
		index2 := rand.Intn(size)
		x1, y1 := fp.Take(index1)
		x2, y2 := fp.Take(index2)
		p1 := c.segmentize(x1, y1, x2, y2)
		p2 := c.segmentizeLinear(x1, y1, x2, y2)
		fmt.Println(x1, y1, x2, y2)
		fmt.Println(p1)
		fmt.Println(p2)
		compareConcaveHulls(t, p1, p2)
	}

}
//go:generate git submodule init
//go:generate git submodule update
//go:generate git submodule foreach git pull origin master
func BenchmarkCompute_ConcaveHullSmall (b * testing.B) {
	dir := "examples/examples"
	err := filepath.Walk(dir, func (path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !info.IsDir() && info.Name() != ".git" {
			b.Run("CPU", func (b *testing.B) {
				scanBenchmark(b, path, info, false)
			})
			b.Run("Memory", func (b *testing.B) {
				scanBenchmark(b, path, info, true)
			})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompute_ConcaveHullSmall (t * testing.T) {
	dir := "examples/examples"
	err := filepath.Walk(dir, func (path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !info.IsDir() && info.Name() != ".git" {
			testExample(t, path, info)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func Benchmark_ConcaveHullBig (b * testing.B) {
	dir := "examples/large-examples"
	err := filepath.Walk(dir, func (path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !info.IsDir() && info.Name() != ".git" {
			scanBenchmark(b, path, info, false)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func scanBenchmark (b * testing.B, path string, f os.FileInfo, isMemoryTest bool) {
	points := readExampleFile(path, f)
	b.Run(path, func (b * testing.B) {
		if isMemoryTest {
			b.ReportAllocs()
		}
		for n := 0; n < b.N; n++ {
			_ = ComputeWithOptions(points, &Options{seglength: 1}) // coordinates are in a projection
		}
	})
}

func testExample (t * testing.T, path string, f os.FileInfo) {
	points := readExampleFile(path, f)
	expected := readResultFile(filepath.Join("examples-result", f.Name()))
	t.Run(path, func (t * testing.T) {
		concaveHull := Compute(points)
		assert.Equal(t, expected, []float64(concaveHull))
	})
}
func readResultFile(path string) []float64 {
	file, err := os.Open(path)
	if (err != nil) {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var points []float64
	for scanner.Scan() {
		myfunc := func(c rune) bool {
			return c == ' '
		}
		coordinates := strings.FieldsFunc(scanner.Text(), myfunc)
		for _, fs := range(coordinates) {
			f, err := strconv.ParseFloat(fs, 64)
			if err != nil {
				log.Fatal(err)
			}
			points = append(points, f)
		}
	}
	return points
}

func readExampleFile(path string, f os.FileInfo) []float64 {
	file, err := os.Open(path)
	if (err != nil) {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var points []float64
	for scanner.Scan() {
		myfunc := func(c rune) bool {
			return c == ';'
		}
		coordinates := strings.FieldsFunc(scanner.Text(), myfunc)
		x, _ := strconv.ParseFloat(coordinates[1], 64)
		y, err := strconv.ParseFloat(coordinates[0], 64)
		// mainly headers
		if (err == nil) {
			points = append(points, x, y)
		}
	}
	return points
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

// Split side in small edges, for each edge find closest point. Remove duplicates
func (c * concaver) segmentizeLinear (x1, y1, x2, y2 float64) (points []float64) {
	dist := math.Sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2))
	nSegments := math.Ceil(dist / c.seglength)
	factor := 1 / nSegments
	flatPoints := make([]float64, 0, int(2 * nSegments))
	vX := factor * (x2 - x1)
	vY := factor * (y2 - y1)

	currentX := x1
	currentY := y1

	latestX := x1
	latestY := y1
	for i := 0; i < int(nSegments); i++ {
		x, y, _, _ := c.rtree.FindNearestPoint(currentX, currentY)
		if x != latestX || y != latestY {
			flatPoints = append(flatPoints, x, y)
			latestX = x
			latestY = y
		}
		currentX += vX
		currentY += vY
	}
	return flatPoints
}

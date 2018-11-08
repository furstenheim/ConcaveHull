## Concave Hull

### Usage

The algorithms accepts a flat array of coordinates (even positions are latitudes and odd positions are longitudes)

    coordinates = []float64{x0, y0, x1, y1, ...}
    concaveHull := ConcaveHull.Compute(ConcaveHull.FlatPoints(coordinates))

### Algorithm

The algorithm starts from a convex hull of the given points and find points close to the edges to build the final polygon. Finally Douglas Peucker is applied to simplify the polygon.
It builds a Concave Hull around the points but it is not an [alpha shape](https://en.wikipedia.org/wiki/Alpha_shape).




### Performance

The following benchmark was run on example 4 from [this website](https://www.codeproject.com/Articles/1201438/The-Concave-Hull-of-a-Set-of-Points). It took 0.19s to build the concave hull. The benchmark was done in a i5 2.50GHz 8Gb of RAM running on Linux

    BenchmarkCompute_ConcaveHullSmall/CPU#03/examples/examples/4-camps-drift.txt-4       	      20	 193270383 ns/op

![Concave hull of a network](./example.png?raw=true "Concave Hull")

To run the benchmarks run `go generate`, this will download all the necessary files

### Installation

This project uses [dep](https://github.com/golang/dep) to handle dependencies.


### Acknowledgments

The algorithm is based on [SnapHull](https://github.com/skipperkongen/jts-algorithm-pack/blob/master/src/org/geodelivery/jap/concavehull/SnapHull.java) from Pimin Konstantin Kefaloukos, which in turn is based on the algorithm for st_concavehull from Postgis 2.0.
Thanks to [acraig](https://github.com/acraig5075) for providing the data set of the picture

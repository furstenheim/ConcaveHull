bench:
	go test -v -bench=ConcaveHullSmall --benchtime=3s
bench-mem:
	go test -v -bench=ConcaveHullSmall --benchtime=3s -benchmem
bench-graph:
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench ConcaveHullSmall -cpuprofile benchmarks/$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg ConcaveHull.test benchmarks/$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(git rev-parse HEAD)/cpu.svg

bench-all:
	go test -v -bench=. --benchtime=3s
bench:
	go test -v -bench=ConcaveHullSmall/CPU --benchtime=3s
bench-mem:
	go test -v -bench=ConcaveHullSmall/Memory --benchtime=3s -benchmem
bench-mem-trace:
	go test -c
	GODEBUG=allocfreetrace=1 ./ConcaveHull.test -test.run=none -test.benchtime=10ms -test.bench=ConcaveHullSmall/Memory#03/examples/examples/4-camps-drift.txt 2>trace.log
bench-graph:
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench ConcaveHullSmall -cpuprofile benchmarks/$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg ConcaveHull.test benchmarks/$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(git rev-parse HEAD)/cpu.svg
bench-mem-alloc-graph:
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench ConcaveHullSmall -memprofile benchmarks/$$(git rev-parse HEAD)/heap.prof
	go tool pprof -lines -sample_index=alloc_objects -gif ConcaveHull.test benchmarks/$$(git rev-parse HEAD)/heap.prof > benchmarks/$$(git rev-parse HEAD)/heap.gif

bench-segmentize:
	go test -v -bench=segmentize --benchtime=3s
bench-segmentize-mem:
	go test -v -bench=segmentize --benchtime=3s -benchmem
bench-segmentize-graph:
	mkdir -p benchmarks/$$(git rev-parse HEAD)
	go test -run=XXX -bench segmentize -cpuprofile benchmarks/$$(git rev-parse HEAD)/cpu.prof
	go tool pprof -svg ConcaveHull.test benchmarks/$$(git rev-parse HEAD)/cpu.prof > benchmarks/$$(git rev-parse HEAD)/cpu.svg
bench-segmentize-mem-trace:
	go test -c
	GODEBUG=allocfreetrace=1 ./ConcaveHull.test -test.run=none -test.benchtime=10ms -test.bench=Benchmark_segmentize/200000 2>trace.log

bench-all:
	go test -v -bench=. --benchtime=3s
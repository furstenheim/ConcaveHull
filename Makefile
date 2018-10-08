bench:
	go test -v -bench=ConcaveHullSmall --benchtime=3s
bench-all:
	go test -v -bench=. --benchtime=3s
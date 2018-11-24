BenchmarkCompute_ConcaveHullSmall/Memory/examples/examples/1-kwanyamazane.txt-4         	   10000	    576615 ns/op	    7937 B/op	      23 allocs/op
BenchmarkCompute_ConcaveHullSmall/Memory#01/examples/examples/2-DT71_045.txt-4          	     500	   8431159 ns/op	   47550 B/op	      26 allocs/op
BenchmarkCompute_ConcaveHullSmall/Memory#02/examples/examples/3-table-mountain.txt-4    	     200	  27901087 ns/op	   69359 B/op	      27 allocs/op
BenchmarkCompute_ConcaveHullSmall/Memory#03/examples/examples/4-camps-drift.txt-4       	      20	 191443688 ns/op	  973547 B/op	      29 allocs/op


BenchmarkCompute_ConcaveHullSmall/CPU/examples/examples/1-kwanyamazane.txt-4         	   10000	    566676 ns/op
BenchmarkCompute_ConcaveHullSmall/CPU#01/examples/examples/2-DT71_045.txt-4          	     500	   8236863 ns/op
BenchmarkCompute_ConcaveHullSmall/CPU#02/examples/examples/3-table-mountain.txt-4    	     200	  31745043 ns/op
BenchmarkCompute_ConcaveHullSmall/CPU#03/examples/examples/4-camps-drift.txt-4       	      20	 179544862 ns/op


## low alloc are due to pooling
Benchmark_segmentize/10-4         	 2000000	      1967 ns/op	      31 B/op	       0 allocs/op
Benchmark_segmentize/1000-4       	  200000	     35615 ns/op	      31 B/op	       0 allocs/op
Benchmark_segmentize/10000-4      	   50000	     97361 ns/op	      32 B/op	       0 allocs/op
Benchmark_segmentize/100000-4     	   20000	    239423 ns/op	      33 B/op	       0 allocs/op
Benchmark_segmentize/200000-4     	   20000	    281185 ns/op	      33 B/op	       0 allocs/op
Benchmark_segmentize/1000000-4    	   10000	    410572 ns/op	      36 B/op	       1 allocs/op
Benchmark_segmentize/10-4         	10000000	       356 ns/op
Benchmark_segmentize/1000-4       	  200000	     20340 ns/op
Benchmark_segmentize/10000-4      	   50000	     79732 ns/op
Benchmark_segmentize/100000-4     	   30000	    182115 ns/op
Benchmark_segmentize/200000-4     	   20000	    210123 ns/op
Benchmark_segmentize/1000000-4    	   10000	    301637 ns/op


Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/100000.txt-4              	     100	  60669943 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/1000000.txt-4             	       5	 630579417 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/200000.txt-4              	      30	 121172644 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/300000.txt-4              	      20	 175455501 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/400000.txt-4              	      20	 240330054 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/500000.txt-4              	      10	 364067480 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/600000.txt-4              	      10	 416798024 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/700000.txt-4              	      10	 446509258 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/800000.txt-4              	      10	 500076069 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/1-zeerust/900000.txt-4              	      10	 563148651 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/100000.txt-4             	     100	  59016604 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/1000000.txt-4            	       5	 657621094 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/200000.txt-4             	      30	 129833518 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/300000.txt-4             	      20	 178284810 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/400000.txt-4             	      20	 233572306 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/500000.txt-4             	      20	 300784692 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/600000.txt-4             	      10	 385247379 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/700000.txt-4             	      10	 459859753 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/800000.txt-4             	      10	 506098878 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/2-overburg/900000.txt-4             	      10	 601230615 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/100000.txt-4            	      50	  76719691 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/1000000.txt-4           	       5	 787979802 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/200000.txt-4            	      30	 142008575 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/300000.txt-4            	      20	 202408719 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/400000.txt-4            	      20	 262201240 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/500000.txt-4            	      10	 323198343 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/600000.txt-4            	      10	 457427329 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/700000.txt-4            	      10	 569549343 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/800000.txt-4            	       5	 608717038 ns/op
Benchmark_ConcaveHullBig/examples/large-examples/3-geometric/900000.txt-4            	       5	 699770391 ns/op

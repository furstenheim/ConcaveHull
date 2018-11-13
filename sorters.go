package ConcaveHull

type lexSorter FlatPoints

func (s lexSorter) Less (i, j int) bool {
	if s[2 * i] < s[2 * j] {
		return true
	}
	if s[2 * i] > s[2 * j] {
		return false
	}

	if s[2 * i + 1] < s[2 * j + 1] {
		return true
	}

	if s[2 * i + 1] > s[2 * j + 1] {
		return false
	}
	return true
}

func (s lexSorter) Len () (int) {
	return len(s) / 2
}

func (s lexSorter) Swap (i, j int) {
	s[2 * i], s[2 * i + 1], s[2 * j], s[2 * j + 1] = s[2 * j], s[2 * j + 1], s[2 * i], s[2 * i + 1]
}

type closestPointSorter []closestPoint

func (s closestPointSorter) Less (i, j int) bool {
	return s[i].index < s[j].index
}

func (s closestPointSorter) Len () (int) {
	return len(s)
}

func (s closestPointSorter) Swap (i, j int) {
	s[i], s[j] = s[j], s[i]
}

// inlined for small boost
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Insertion sort
func (c closestPointSorter) cpInsertionSort(a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && c.Less(j, j-1); j-- {
			c.Swap(j, j-1)
		}
	}
}

// siftDown implements the heap property on c[lo, hi).
// first is an offset into the array where the root of the heap lies.
func (c closestPointSorter) cpSiftDown(lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && c.Less(first+child, first+child+1) {
			child++
		}
		if !c.Less(first+root, first+child) {
			return
		}
		c.Swap(first+root, first+child)
		root = child
	}
}

func (c closestPointSorter) cpHeapSort(a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		c.cpSiftDown(i, hi, first)
	}

	// Pop elements, largest first, into end of c.
	for i := hi - 1; i >= 0; i-- {
		c.Swap(first, first+i)
		c.cpSiftDown(lo, i, first)
	}
}

// Quicksort, loosely following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThree moves the median of the three values c[m0], c[m1], c[m2] into c[m1].
func (c closestPointSorter) cpMedianOfThree(m1, m0, m2 int) {
	// sort 3 elements
	if c.Less(m1, m0) {
		c.Swap(m1, m0)
	}
	// c[m0] <= c[m1]
	if c.Less(m2, m1) {
		c.Swap(m2, m1)
		// c[m0] <= c[m2] && c[m1] < c[m2]
		if c.Less(m1, m0) {
			c.Swap(m1, m0)
		}
	}
	// now c[m0] <= c[m1] <= c[m2]
}

func (c closestPointSorter) cpSwapRange(a, b, n int) {
	for i := 0; i < n; i++ {
		c.Swap(a+i, b+i)
	}
}

func (cp closestPointSorter) cpDoPivot(lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		cp.cpMedianOfThree(lo, lo+s, lo+2*s)
		cp.cpMedianOfThree(m, m-s, m+s)
		cp.cpMedianOfThree(hi-1, hi-1-s, hi-1-2*s)
	}
	cp.cpMedianOfThree(lo, m, hi-1)

	// Invariants are:
	//	c[lo] = pivot (set up by ChoosePivot)
	//	c[lo < i < a] < pivot
	//	c[a <= i < b] <= pivot
	//	c[b <= i < c] unexamined
	//	c[c <= i < hi-1] > pivot
	//	c[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && cp.Less(a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !cp.Less(pivot, b); b++ { // c[b] <= pivot
		}
		for ; b < c && cp.Less(pivot, c-1); c-- { // c[c-1] > pivot
		}
		if b >= c {
			break
		}
		// c[b] > pivot; c[c-1] <= pivot
		cp.Swap(b, c-1)
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !cp.Less(pivot, hi-1) { // c[hi-1] = pivot
			cp.Swap(c, hi-1)
			c++
			dups++
		}
		if !cp.Less(b-1, pivot) { // c[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> c[m] <= pivot
		if !cp.Less(m, pivot) { // c[m] = pivot
			cp.Swap(m, b-1)
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	c[a <= i < b] unexamined
		//	c[b <= i < c] = pivot
		for {
			for ; a < b && !cp.Less(b-1, pivot); b-- { // c[b] == pivot
			}
			for ; a < b && cp.Less(a, pivot); a++ { // c[a] < pivot
			}
			if a >= b {
				break
			}
			// c[a] == pivot; c[b-1] < pivot
			cp.Swap(a, b-1)
			a++
			b--
		}
	}
	// Swap pivot into middle
	cp.Swap(pivot, b-1)
	return b - 1, c
}

func (c closestPointSorter) cpQuickSort(a, b, maxDepth int) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			c.cpHeapSort(a, b)
			return
		}
		maxDepth--
		mlo, mhi := c.cpDoPivot(a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			c.cpQuickSort(a, mlo, maxDepth)
			a = mhi // i.e., (c closestPointSorter) c.cpQuickSort(mhi, b)
		} else {
			c.cpQuickSort(mhi, b, maxDepth)
			b = mlo // i.e., (c closestPointSorter) cpQuickSort(c, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if c.Less(i, i-6) {
				c.Swap(i, i-6)
			}
		}
		c.cpInsertionSort(a, b)
	}
}

// Sort sorts c.
// It makes one call to c.Len to determine n, and O(n*log(n)) calls to
// c.Less and c.Swap. The sort is not guaranteed to be stable.
func (c closestPointSorter) cpSort () {
	n := c.Len()
	c.cpQuickSort(0, n, maxDepth(n))
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returns 2*ceil(lg(n+1)).
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

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

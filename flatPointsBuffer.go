package ConcaveHull

type flatPointsBuffer struct {
	arrays [][]float64
	currentCap int
	currentIndex int
}

func makeFlatPointBuffer (size int) flatPointsBuffer {
	arrays := make([][]float64, 0, 5)
	firstArray := make([]float64, 0, size)
	arrays = append(arrays, firstArray)
	fpb := flatPointsBuffer{
		arrays: arrays,
		currentIndex: 0,
		currentCap: 1,
	}
	return fpb
}

func (fpb *flatPointsBuffer) addFloat (x float64) {
	currentArray := fpb.arrays[fpb.currentIndex]
	if len(currentArray) == cap(currentArray) {
		fpb.currentIndex++
		if fpb.currentIndex >= fpb.currentCap {
			fpb.arrays = append(fpb.arrays, make([]float64, 0, 2 * cap(currentArray)))
		}
	}
	fpb.arrays[fpb.currentIndex] = append(fpb.arrays[fpb.currentIndex], x)
}
func (fpb *flatPointsBuffer) toFloatArray () []float64 {
	totalLength := 0
	for i := 0; i <= fpb.currentIndex; i++ {
		totalLength += len(fpb.arrays[i])
	}
	result := make([]float64, 0, totalLength)

	for i := 0; i <= fpb.currentIndex; i++ {
		result = append(result, fpb.arrays[i]...)
	}
	return result
}



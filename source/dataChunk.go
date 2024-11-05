package main

// Primitive data struct to construct bigger chunks
type dataChunk struct {
	data [][]*swimmingAnimal
	// Pointers to surrounding chunks
	aboveDataChunk *dataChunk
	belowDataChunk *dataChunk
}

// Right index with wrap around
func (dC dataChunk) getRightIndex(index int) int {
	return mod(index+1, len(dC.data[0]))
}

// Left index with wrap around
func (dC dataChunk) getLeftIndex(index int) int {
	return mod(index-1, len(dC.data[0]))
}

// Get above animal which may be in a different chunk
func (dC dataChunk) getAboveAnimal(indexX, indexY int) *swimmingAnimal {
	if indexY - 1 < 0 {
		return dC.aboveDataChunk.data[len(dC.aboveDataChunk.data)-1][indexX]
	} else {
		return dC.data[indexY-1][indexX]
	}
}

func (dC dataChunk) getBellowAnimal(indexX, indexY int) *swimmingAnimal {
	return dC.data[indexY][indexX]
}

func (dC dataChunk) process() {
	for
}

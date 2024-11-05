package main

// Primitive data struct to construct bigger chunks
type dataChunk struct {
	data [][]*swimmingAnimal
	// Pointers to surrounding chunks
	aboveDataChunk *dataChunk
	belowDataChunk *dataChunk
}

// Right index with wrap around
func (dC dataChunk) getRightIndex(index, gridSizeX int) int {
	return mod(index+1, gridSizeX)
}

// Left index with wrap around
func (dC dataChunk) getLeftIndex(index, gridSizeX int) int {
	return mod(index-1, gridSizeX)
}

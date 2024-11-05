package main

// Thread chunk is used for holding data for each thread
type threadChunk struct {
	data [][]*swimmingAnimal
	// Pointers to surrounding chunks
	aboveBorderChunk *borderChunk
	belowBorderChunk *borderChunk
}

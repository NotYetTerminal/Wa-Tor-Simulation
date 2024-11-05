package main

// Modulo for both positive and negative numbers
func mod(a, b int) int {
	return (a%b + b) % b
}

// Converts absolute index to X and Y positions
func indexToPosition(index, gridSizeX int) (posX, posY int) {
	return mod(index, gridSizeX), index / gridSizeX
}

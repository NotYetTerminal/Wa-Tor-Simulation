// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-05
// Description:
// Common methods used by other files

package main

// Modulo for both positive and negative numbers
func mod(a, b int) int {
	return (a%b + b) % b
}

// Converts absolute index to X and Y positions
func indexToPosition(index int) (posX, posY int) {
	return mod(index, GridSizeX), index / GridSizeX
}

// Left index with wrap around
func getLeftIndex(index int) int {
	return mod(index-1, GridSizeX)
}

// Right index with wrap around
func getRightIndex(index int) int {
	return mod(index+1, GridSizeX)
}

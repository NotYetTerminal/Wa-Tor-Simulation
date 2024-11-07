// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-07
// Description:
// File for storing code related to rendering

package main

import (
	"fmt"
	"image/color"
)

// Updates the graphical screen
func updateScreen(tC *threadChunk) {
	currentTC := tC
	var allData [][]color.RGBA
	for range Threads {
		// Thread data
		for _, row := range currentTC.data {
			var rowColourSlice []color.RGBA
			for _, animal := range row {
				rowColourSlice = append(rowColourSlice, convertAnimalToColour(animal))
			}
			allData = append(allData, rowColourSlice)
		}
		allData = append(allData, []color.RGBA{})
		// Border data
		for _, row := range currentTC.belowBorderChunk.data {
			var rowColourSlice []color.RGBA
			for _, animal := range row {
				rowColourSlice = append(rowColourSlice, convertAnimalToColour(animal))
			}
			allData = append(allData, rowColourSlice)
		}
		allData = append(allData, []color.RGBA{})
		currentTC = currentTC.belowBorderChunk.belowThreadChunk
	}
	for _, colour := range allData {
		fmt.Println(colour)
	}
}

// Converts swimmingAnimal to pixel colours
func convertAnimalToColour(animal *swimmingAnimal) color.RGBA {
	if animal == nil {
		return color.RGBA{A: 255} // Black
	} else if animal.isShark {
		return color.RGBA{R: 255, A: 255} // Red
	} else if animal.isShark {
		return color.RGBA{B: 255, A: 255} // Blue
	} else {
		return color.RGBA{A: 255} // Black
	}
}

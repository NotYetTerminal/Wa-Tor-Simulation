// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-07
// Description:
// File for storing code related to rendering

package main

import (
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"time"
)

// Updates the graphical screen
func updateScreen(tC *threadChunk) {
	currentTC := tC
	var absoluteY int
	for range Threads {
		// Thread data
		for indexY, row := range currentTC.data {
			for indexX, animal := range row {
				rgbaImage.Set(indexX, absoluteY+indexY, convertAnimalToColour(animal))
			}
		}
		absoluteY += len(currentTC.data)
		// Border data
		for indexY, row := range currentTC.belowBorderChunk.data {
			for indexX, animal := range row {
				rgbaImage.Set(indexX, absoluteY+indexY, convertAnimalToColour(animal))
			}
		}
		absoluteY += len(currentTC.belowBorderChunk.data)
		currentTC = currentTC.belowBorderChunk.belowThreadChunk
	}
	canvas.Refresh(canvasToWrite)
	time.Sleep(time.Duration(1) * time.Millisecond)
}

// Converts swimmingAnimal to pixel colours
func convertAnimalToColour(animal *swimmingAnimal) color.RGBA {
	if animal == nil {
		return color.RGBA{A: 255} // Black
	} else if animal.isShark {
		return color.RGBA{R: 255, A: 255} // Red
	} else if !animal.isShark {
		return color.RGBA{B: 255, A: 255} // Blue
	} else {
		return color.RGBA{A: 255} // Black
	}
}

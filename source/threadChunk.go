// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-05
// Description:
// File for storing code related to the thread chunks

package main

// Thread chunk is used for holding data for each thread
// Pointers are to surrounding chunks
type threadChunk struct {
	data             [][]*swimmingAnimal
	aboveBorderChunk *borderChunk
	belowBorderChunk *borderChunk
}

// Get one row
func (tC threadChunk) getRow(indexY int) []*swimmingAnimal {
	return tC.data[indexY]
}

// Getters for swimmingAnimal
// May be gotten from a different chunk
func (tC threadChunk) getAnimal(indexX, indexY int) *swimmingAnimal {
	return tC.data[indexY][indexX]
}

func (tC threadChunk) getLeftAnimal(indexX, indexY int) *swimmingAnimal {
	return tC.data[indexY][getLeftIndex(indexX)]
}

func (tC threadChunk) getRightAnimal(indexX, indexY int) *swimmingAnimal {
	return tC.data[indexY][getRightIndex(indexX)]
}

func (tC threadChunk) getAboveAnimal(indexX, indexY int) *swimmingAnimal {
	if indexY-1 < 0 {
		return tC.aboveBorderChunk.getBottomRowAnimal(indexX)
	} else {
		return tC.data[indexY-1][indexX]
	}
}

func (tC threadChunk) getBelowAnimal(indexX, indexY int) *swimmingAnimal {
	if indexY+1 >= len(tC.data) {
		return tC.belowBorderChunk.getTopRowAnimal(indexX)
	} else {
		return tC.data[indexY+1][indexX]
	}
}

// Setters for swimmingAnimal
// May be set in a different chunk
func (tC threadChunk) setAnimal(indexX, indexY int, animal *swimmingAnimal) {
	tC.data[indexY][indexX] = animal
}

func (tC threadChunk) setLeftAnimal(indexX, indexY int, animal *swimmingAnimal) {
	tC.data[indexY][getLeftIndex(indexX)] = animal
}

func (tC threadChunk) setRightAnimal(indexX, indexY int, animal *swimmingAnimal) {
	tC.data[indexY][getRightIndex(indexX)] = animal
}

func (tC threadChunk) setAboveAnimal(indexX, indexY int, animal *swimmingAnimal) {
	if indexY-1 < 0 {
		tC.aboveBorderChunk.setBottomRowAnimal(indexX, animal)
	} else {
		tC.data[indexY-1][indexX] = animal
	}
}

func (tC threadChunk) setBelowAnimal(indexX, indexY int, animal *swimmingAnimal) {
	if indexY+1 >= len(tC.data) {
		tC.belowBorderChunk.setTopRowAnimal(indexX, animal)
	} else {
		tC.data[indexY+1][indexX] = animal
	}
}

// Getters for 2 rows in data
func (tC threadChunk) getTopRowAnimal(indexX int) *swimmingAnimal {
	return tC.data[0][indexX]
}

func (tC threadChunk) getBottomRowAnimal(indexX int) *swimmingAnimal {
	return tC.data[len(tC.data)-1][indexX]
}

// Setters for 2 rows in data
func (tC threadChunk) setTopRowAnimal(indexX int, animal *swimmingAnimal) {
	tC.data[0][indexX] = animal
}

func (tC threadChunk) setBottomRowAnimal(indexX int, animal *swimmingAnimal) {
	tC.data[len(tC.data)-1][indexX] = animal
}

// Process all of the rows managed by thread
func (tC threadChunk) processAllRows(checkingShark bool) {
	// Process above border bottom row and top thread row
	tC.aboveBorderChunk.lockBorderChunk()
	tC.aboveBorderChunk.processBottomBorderRow(checkingShark)
	tC.processTopThreadRow(checkingShark)
	tC.aboveBorderChunk.unlockBorderChunk()

	if len(tC.data) > 2 {
		// Process middle rows
		tC.processMiddleThreadRow(checkingShark)
	}

	// Process below border top row and bottom thread row
	tC.belowBorderChunk.lockBorderChunk()
	tC.processBottomThreadRow(checkingShark)
	tC.belowBorderChunk.processTopBorderRow(checkingShark)
	tC.belowBorderChunk.unlockBorderChunk()
}

// Processes the top row
func (tC threadChunk) processTopThreadRow(checkingShark bool) {
	processRow(checkingShark, 0, &tC)
}

// Processes the middle rows
func (tC threadChunk) processMiddleThreadRow(checkingShark bool) {
	for indexY := range len(tC.data) - 2 {
		processRow(checkingShark, indexY+1, &tC)
	}
}

// Processes the bottom row
func (tC threadChunk) processBottomThreadRow(checkingShark bool) {
	processRow(checkingShark, len(tC.data)-1, &tC)
}

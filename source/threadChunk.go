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

	// Process middle rows
	tC.processMiddleThreadRow(checkingShark)

	// Process below border top row and bottom thread row
	tC.belowBorderChunk.lockBorderChunk()
	tC.processBottomThreadRow(checkingShark)
	tC.belowBorderChunk.processTopBorderRow(checkingShark)
	tC.belowBorderChunk.unlockBorderChunk()
}

// Processes the top row
func (tC threadChunk) processTopThreadRow(checkingShark bool) {
	tC.processThreadRow(checkingShark, 0)
}

// Processes the middle rows
func (tC threadChunk) processMiddleThreadRow(checkingShark bool) {
	for indexY := range tC.data[1 : len(tC.data)-1] {
		tC.processThreadRow(checkingShark, indexY)
	}
}

// Processes the bottom row
func (tC threadChunk) processBottomThreadRow(checkingShark bool) {
	tC.processThreadRow(checkingShark, len(tC.data)-1)
}

// Runs through a single row and processes all animals signaled
func (tC threadChunk) processThreadRow(checkingShark bool, indexY int) {
	for indexX, animal := range tC.getRow(indexY) {
		// If checking states are the same it means that this animal has been already checked this iteration
		if animal != nil && animal.checkingState != CurrentCheckingState {
			// Check sharks
			if checkingShark && animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get whether fish are nearby and valid movement directions
				movingToFish, direction := getFishOrFreeSpace(tC.getLeftAnimal(indexX, indexY), tC.getRightAnimal(indexX, indexY), tC.getAboveAnimal(indexX, indexY), tC.getBelowAnimal(indexX, indexY))

				// Shark eating fish while moving
				if movingToFish {
					animal.energy = Starve
					fishCount.Add(-1)
				} else {
					animal.energy -= 1
					// Shark died
					if animal.energy == 0 {
						tC.moveAnimal(nil, nil, indexX, indexY, 0)
						sharkCount.Add(-1)
						continue
					}
				}

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					if direction != 0 {
						tC.moveAnimal(animal, newSwimmingAnimal(true, CurrentCheckingState, SharkBreed, Starve), indexX, indexY, direction)
						sharkCount.Add(1)
					}
					animal.reproductionTime = SharkBreed
					continue
				}

				// Move to square
				tC.moveAnimal(animal, nil, indexX, indexY, direction)
			} else if !checkingShark && !animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get move direction
				direction := getFreeSpace(tC.getLeftAnimal(indexX, indexY), tC.getRightAnimal(indexX, indexY), tC.getAboveAnimal(indexX, indexY), tC.getBelowAnimal(indexX, indexY))

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					if direction != 0 {
						tC.moveAnimal(animal, newSwimmingAnimal(false, CurrentCheckingState, FishBreed, 0), indexX, indexY, direction)
						fishCount.Add(1)
					}
					animal.reproductionTime = FishBreed
					continue
				}

				// Move to square
				tC.moveAnimal(animal, nil, indexX, indexY, direction)
			}
		}
	}
}

// Modify moving to and moving from square in dataChunk
func (tC threadChunk) moveAnimal(movingToSquare, movingFromSquare *swimmingAnimal, indexX, indexY, direction int) {
	switch direction {
	case 1:
		tC.setLeftAnimal(indexX, indexY, movingToSquare)
	case 2:
		tC.setRightAnimal(indexX, indexY, movingToSquare)
	case 3:
		tC.setAboveAnimal(indexX, indexY, movingToSquare)
	case 4:
		tC.setBelowAnimal(indexX, indexY, movingToSquare)
	}
	tC.setAnimal(indexX, indexY, movingFromSquare)
}

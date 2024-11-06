// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-05
// Description:
// File for storing code related to the border chunks

package main

import (
	"sync"
)

// Border chunk is used for holding data between thread chunks
// Pointers are to surrounding chunks
// Mutex used to regulate 2 threads accessing this border chunk
type borderChunk struct {
	data             [][]*swimmingAnimal
	aboveThreadChunk *threadChunk
	belowThreadChunk *threadChunk
	lock             sync.Mutex
}

// Constructor
func newBorderChunk(topData, bottomData []*swimmingAnimal, aboveThreadChunk, belowThreadChunk *threadChunk) *borderChunk {
	return &borderChunk{data: [][]*swimmingAnimal{topData, bottomData}, aboveThreadChunk: aboveThreadChunk, belowThreadChunk: belowThreadChunk, lock: sync.Mutex{}}
}

// Mutex locking
func (bC *borderChunk) lockBorderChunk() {
	bC.lock.Lock()
}

func (bC *borderChunk) unlockBorderChunk() {
	bC.lock.Unlock()
}

// Get one row
func (bC *borderChunk) getRow(indexY int) []*swimmingAnimal {
	return bC.data[indexY]
}

// Getters for swimmingAnimal
// May be set in a different chunk
func (bC *borderChunk) getAnimal(indexX, indexY int) *swimmingAnimal {
	return bC.data[indexY][indexX]
}

func (bC *borderChunk) getLeftAnimal(indexX, indexY int) *swimmingAnimal {
	return bC.data[indexY][getLeftIndex(indexX)]
}

func (bC *borderChunk) getRightAnimal(indexX, indexY int) *swimmingAnimal {
	return bC.data[indexY][getRightIndex(indexX)]
}

func (bC *borderChunk) getAboveAnimal(indexX, indexY int) *swimmingAnimal {
	if indexY-1 < 0 {
		return bC.aboveThreadChunk.getBottomRowAnimal(indexX)
	} else {
		return bC.data[indexY-1][indexX]
	}
}

func (bC *borderChunk) getBelowAnimal(indexX, indexY int) *swimmingAnimal {
	if indexY+1 >= len(bC.data) {
		return bC.belowThreadChunk.getTopRowAnimal(indexX)
	} else {
		return bC.data[indexY+1][indexX]
	}
}

// Setters for swimmingAnimal
// May be set in a different chunk
func (bC *borderChunk) setAnimal(indexX, indexY int, animal *swimmingAnimal) {
	bC.data[indexY][indexX] = animal
}

func (bC *borderChunk) setLeftAnimal(indexX, indexY int, animal *swimmingAnimal) {
	bC.data[indexY][getLeftIndex(indexX)] = animal
}

func (bC *borderChunk) setRightAnimal(indexX, indexY int, animal *swimmingAnimal) {
	bC.data[indexY][getRightIndex(indexX)] = animal
}

func (bC *borderChunk) setAboveAnimal(indexX, indexY int, animal *swimmingAnimal) {
	if indexY-1 < 0 {
		bC.aboveThreadChunk.setBottomRowAnimal(indexX, animal)
	} else {
		bC.data[indexY-1][indexX] = animal
	}
}

func (bC *borderChunk) setBelowAnimal(indexX, indexY int, animal *swimmingAnimal) {
	if indexY+1 >= len(bC.data) {
		bC.belowThreadChunk.setTopRowAnimal(indexX, animal)
	} else {
		bC.data[indexY+1][indexX] = animal
	}
}

// Getters for 2 rows in data
func (bC *borderChunk) getTopRowAnimal(indexX int) *swimmingAnimal {
	return bC.data[0][indexX]
}

func (bC *borderChunk) getBottomRowAnimal(indexX int) *swimmingAnimal {
	return bC.data[1][indexX]
}

// Setters for 2 rows in data
func (bC *borderChunk) setTopRowAnimal(indexX int, animal *swimmingAnimal) {
	bC.data[0][indexX] = animal
}

func (bC *borderChunk) setBottomRowAnimal(indexX int, animal *swimmingAnimal) {
	bC.data[1][indexX] = animal
}

// Used by the thread above the border
func (bC *borderChunk) processTopBorderRow(checkingShark bool) {
	bC.processBorderRow(checkingShark, 0)
}

// Used by the thread below the border
func (bC *borderChunk) processBottomBorderRow(checkingShark bool) {
	bC.processBorderRow(checkingShark, 1)
}

// Runs through a single row and processes all animals signaled
func (bC *borderChunk) processBorderRow(checkingShark bool, indexY int) {
	for indexX, animal := range bC.getRow(indexY) {
		// If checking states are the same it means that this animal has been already checked this iteration
		if animal != nil && animal.checkingState != CurrentCheckingState {
			// Check sharks
			if checkingShark && animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get whether fish are nearby and valid movement directions
				movingToFish, direction := getFishOrFreeSpace(bC.getLeftAnimal(indexX, indexY), bC.getRightAnimal(indexX, indexY), bC.getAboveAnimal(indexX, indexY), bC.getBelowAnimal(indexX, indexY))

				// Shark eating fish while moving
				if movingToFish {
					animal.energy = Starve
				} else {
					animal.energy -= 1
					// Shark died
					if animal.energy == 0 {
						bC.moveAnimal(nil, nil, indexX, indexY, 0)
						continue
					}
				}

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					bC.moveAnimal(animal, newSwimmingAnimal(true, CurrentCheckingState, SharkBreed, Starve), indexX, indexY, direction)
					animal.reproductionTime = SharkBreed
					continue
				}

				// Move to square
				bC.moveAnimal(animal, nil, indexX, indexY, direction)
			} else if !checkingShark && !animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get move direction
				direction := getFreeSpace(bC.getLeftAnimal(indexX, indexY), bC.getRightAnimal(indexX, indexY), bC.getAboveAnimal(indexX, indexY), bC.getBelowAnimal(indexX, indexY))

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					bC.moveAnimal(animal, newSwimmingAnimal(false, CurrentCheckingState, FishBreed, 0), indexX, indexY, direction)
					animal.reproductionTime = FishBreed
					continue
				}

				// Move to square
				bC.moveAnimal(animal, nil, indexX, indexY, direction)
			}
		}
	}
}

// Modify moving to and moving from square in dataChunk
func (bC *borderChunk) moveAnimal(movingToSquare, movingFromSquare *swimmingAnimal, indexX, indexY, direction int) {
	switch direction {
	case 1:
		bC.setLeftAnimal(indexX, indexY, movingToSquare)
	case 2:
		bC.setRightAnimal(indexX, indexY, movingToSquare)
	case 3:
		bC.setAboveAnimal(indexX, indexY, movingToSquare)
	case 4:
		bC.setBelowAnimal(indexX, indexY, movingToSquare)
	}
	bC.setAnimal(indexX, indexY, movingFromSquare)
}

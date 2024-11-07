// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-06
// Description:
// File for storing code related to processing a row of animals

package main

import (
	"math/rand/v2"
)

// Runs through a single row and processes all animals signaled
func processRow[T dataChunk](checkingShark bool, indexY int, dC *T) {
	for indexX, animal := range (*dC).getRow(indexY) {
		// If checking states are the same it means that this animal has been already checked this iteration
		if animal != nil && animal.checkingState != CurrentCheckingState {
			// Check sharks
			if checkingShark && animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get whether fish are nearby and valid movement directions
				movingToFish, direction := getFishOrFreeSpace((*dC).getLeftAnimal(indexX, indexY), (*dC).getRightAnimal(indexX, indexY), (*dC).getAboveAnimal(indexX, indexY), (*dC).getBelowAnimal(indexX, indexY))

				// Shark eating fish while moving
				if movingToFish {
					animal.energy = Starve
					fishCount.Add(-1)
				} else {
					animal.energy -= 1
					// Shark died
					if animal.energy == 0 {
						moveAnimal(nil, nil, indexX, indexY, 0, dC)
						sharkCount.Add(-1)
						continue
					}
				}

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					if direction != 0 {
						moveAnimal(animal, newSwimmingAnimal(true, CurrentCheckingState, SharkBreed, Starve), indexX, indexY, direction, dC)
						sharkCount.Add(1)
					}
					animal.reproductionTime = SharkBreed
					continue
				}

				// Move to square
				moveAnimal(animal, nil, indexX, indexY, direction, dC)
			} else if !checkingShark && !animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get move direction
				direction := getFreeSpace((*dC).getLeftAnimal(indexX, indexY), (*dC).getRightAnimal(indexX, indexY), (*dC).getAboveAnimal(indexX, indexY), (*dC).getBelowAnimal(indexX, indexY))

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					if direction != 0 {
						moveAnimal(animal, newSwimmingAnimal(false, CurrentCheckingState, FishBreed, 0), indexX, indexY, direction, dC)
						fishCount.Add(1)
					}
					animal.reproductionTime = FishBreed
					continue
				}

				// Move to square
				moveAnimal(animal, nil, indexX, indexY, direction, dC)
			}
		}
	}
}

// Modify moving to and moving from square in dataChunk
func moveAnimal[T dataChunk](movingToSquare, movingFromSquare *swimmingAnimal, indexX, indexY, direction int, dC *T) {
	switch direction {
	case 1:
		(*dC).setLeftAnimal(indexX, indexY, movingToSquare)
	case 2:
		(*dC).setRightAnimal(indexX, indexY, movingToSquare)
	case 3:
		(*dC).setAboveAnimal(indexX, indexY, movingToSquare)
	case 4:
		(*dC).setBelowAnimal(indexX, indexY, movingToSquare)
	}
	(*dC).setAnimal(indexX, indexY, movingFromSquare)
}

// Gets a list of directions for fish or free spaces
func getFishOrFreeSpace(leftAnimal, rightAnimal, aboveAnimal, belowAnimal *swimmingAnimal) (bool, int) {
	movingToFish := false
	var fishSlice []int
	var freeSpacesSlice []int

	checkAnimal(leftAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 1)
	checkAnimal(rightAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 2)
	checkAnimal(aboveAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 3)
	checkAnimal(belowAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 4)

	if movingToFish {
		return movingToFish, fishSlice[rand.IntN(len(fishSlice))]
	} else if len(freeSpacesSlice) != 0 {
		return movingToFish, freeSpacesSlice[rand.IntN(len(freeSpacesSlice))]
	} else {
		return movingToFish, 0
	}
}

// Checks whether a space is empty or a fish
func checkAnimal(animal *swimmingAnimal, movingToFish *bool, fishSlice, freeSpacesSlice *[]int, direction int) {
	if animal == nil {
		*freeSpacesSlice = append(*freeSpacesSlice, direction)
	} else if !animal.isShark {
		*fishSlice = append(*fishSlice, direction)
		*movingToFish = true
	}
}

// Get free space next to animal
func getFreeSpace(leftAnimal, rightAnimal, aboveAnimal, belowAnimal *swimmingAnimal) int {
	var freeSpacesSlice []int

	if leftAnimal == nil {
		freeSpacesSlice = append(freeSpacesSlice, 1)
	}
	if rightAnimal == nil {
		freeSpacesSlice = append(freeSpacesSlice, 2)
	}
	if aboveAnimal == nil {
		freeSpacesSlice = append(freeSpacesSlice, 3)
	}
	if belowAnimal == nil {
		freeSpacesSlice = append(freeSpacesSlice, 4)
	}

	if len(freeSpacesSlice) != 0 {
		return freeSpacesSlice[rand.IntN(len(freeSpacesSlice))]
	} else {
		return 0
	}
}

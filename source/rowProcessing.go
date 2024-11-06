// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-06
// Description:
// File for storing code related to processing a row of animals

package main

import "math/rand/v2"

// Runs through a single row and processes all animals signaled
func processRow(checkingShark bool, indexY int, chunk *dataChunk) {
	for indexX, animal := range (*chunk).getRow(indexY) {
		// If checking states are the same it means that this animal has been already checked this iteration
		if animal != nil && animal.checkingState != CurrentCheckingState {
			// Check sharks
			if checkingShark && animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get whether fish are nearby and valid movement directions
				movingToFish, direction := getFishOrFreeSpace((*chunk).getLeftAnimal(indexX, indexY), (*chunk).getRightAnimal(indexX, indexY), (*chunk).getAboveAnimal(indexX, indexY), (*chunk).getBellowAnimal(indexX, indexY))

				// Shark eating fish while moving
				if movingToFish {
					animal.energy = Starve
				} else {
					animal.energy -= 1
					// Shark died
					if animal.energy == 0 {
						moveAnimal(nil, nil, chunk, indexX, indexY, 0)
						continue
					}
				}

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					moveAnimal(animal, newSwimmingAnimal(true, CurrentCheckingState, SharkBreed, Starve), chunk, indexX, indexY, direction)
					animal.reproductionTime = SharkBreed
					continue
				}

				// Move to square
				moveAnimal(animal, nil, chunk, indexX, indexY, direction)
			} else if !checkingShark && !animal.isShark {
				animal.checkingState = CurrentCheckingState
				// Get move direction
				direction := getFreeSpace((*chunk).getLeftAnimal(indexX, indexY), (*chunk).getRightAnimal(indexX, indexY), (*chunk).getAboveAnimal(indexX, indexY), (*chunk).getBellowAnimal(indexX, indexY))

				// Move and reproduce
				animal.reproductionTime -= 1
				if animal.reproductionTime == 0 {
					moveAnimal(animal, newSwimmingAnimal(false, CurrentCheckingState, FishBreed, 0), chunk, indexX, indexY, direction)
					animal.reproductionTime = FishBreed
					continue
				}

				// Move to square
				moveAnimal(animal, nil, chunk, indexX, indexY, direction)
			}
		}
	}
}

// Modify moving to and moving from square in dataChunk
func moveAnimal(movingToSquare, movingFromSquare *swimmingAnimal, chunk *dataChunk, indexX, indexY, direction int) {
	switch direction {
	case 1:
		(*chunk).setLeftAnimal(indexX, indexY, movingToSquare)
	case 2:
		(*chunk).setRightAnimal(indexX, indexY, movingToSquare)
	case 3:
		(*chunk).setAboveAnimal(indexX, indexY, movingToSquare)
	case 4:
		(*chunk).setBellowAnimal(indexX, indexY, movingToSquare)
	}
	(*chunk).setAnimal(indexX, indexY, movingFromSquare)
}

// Gets a list of directions for fish or free spaces
func getFishOrFreeSpace(leftAnimal, rightAnimal, aboveAnimal, bellowAnimal *swimmingAnimal) (bool, int) {
	movingToFish := false
	var fishSlice []int
	var freeSpacesSlice []int

	checkAnimal(leftAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 1)
	checkAnimal(rightAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 2)
	checkAnimal(aboveAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 3)
	checkAnimal(bellowAnimal, &movingToFish, &fishSlice, &freeSpacesSlice, 4)

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
func getFreeSpace(leftAnimal, rightAnimal, aboveAnimal, bellowAnimal *swimmingAnimal) int {
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
	if bellowAnimal == nil {
		freeSpacesSlice = append(freeSpacesSlice, 4)
	}

	if len(freeSpacesSlice) != 0 {
		return freeSpacesSlice[rand.IntN(len(freeSpacesSlice))]
	} else {
		return 0
	}
}

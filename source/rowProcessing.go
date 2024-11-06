// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-06
// Description:
// File for storing code related to processing a row of animals

package main

import "math/rand/v2"

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

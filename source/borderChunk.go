package main

import "sync"

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

// Getters for 2 rows in data
func (bC *borderChunk) getTopRowAnimal(index int) *swimmingAnimal {
	return bC.data[0][index]
}

func (bC *borderChunk) getBottomRowAnimal(index int) *swimmingAnimal {
	return bC.data[1][index]
}

// Setters for 2 rows in data
func (bC *borderChunk) setTopRowAnimal(index int, animal *swimmingAnimal) {
	bC.data[0][index] = animal
}

func (bC *borderChunk) setBottomRowAnimal(index int, animal *swimmingAnimal) {
	bC.data[1][index] = animal
}

// Used by the thread above the border
func (bC *borderChunk) processTopRow(checkingShark bool) {
	bC.lock.Lock()
}

// Used by the thread below the border
func (bC *borderChunk) processBottomRow(checkingShark bool) {
	bC.lock.Lock()
	for index, animal := range bC.data[1] {
		if animal != nil {
			if checkingShark && animal.isShark {

			} else if !checkingShark && !animal.isShark {

			}
		}
	}
}

// Runs through a single row and processes all animals signaled
func processRow(checkingShark bool, data, aboveRow, belowRow []*swimmingAnimal) {
	for index, animal := range data {
		if animal != nil {
			if checkingShark && animal.isShark {
				leftAnimal := data[getLeftIndex(index, len(data))]
				rightAnimal := data[getRightIndex(index, len(data))]
				aboveAnimal := aboveRow[index]
				bellowAnimal := belowRow[index]

			} else if !checkingShark && !animal.isShark {

			}
		}
	}
}

func createDirectionSlices() {}

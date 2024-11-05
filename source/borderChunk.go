package main

import "sync"

// Border chunk is used for passing data between threads
type borderChunk struct {
	topDataChunk    dataChunk
	bottomDataChunk dataChunk
	lock            sync.Mutex
}

// Constructor
func newBorderChunk(topThreadDataChunk, bottomThreadDataChunk *dataChunk) *borderChunk {
	topDataChunk := dataChunk{} // NEED DATA
	bottomDataChunk := dataChunk{}

	topDataChunk.aboveDataChunk = topThreadDataChunk
	topDataChunk.belowDataChunk = &bottomDataChunk

	bottomDataChunk.aboveDataChunk = &topDataChunk
	bottomDataChunk.belowDataChunk = bottomThreadDataChunk

	return &borderChunk{bottomDataChunk: bottomDataChunk, lock: sync.Mutex{}, topDataChunk: topDataChunk}
}

// Process the data before array
func (bC *borderChunk) processDataBefore(checkingShark bool) {
	bC.lock.Lock()
	//for index := range cB.dataBefore {
	//	swimmingAnimal := cB.dataBefore[index]
	//	if swimmingAnimal != nil {
	//		if checkingShark && swimmingAnimal.isShark {
	//			swimmingAnimal.process()
	//		} else if !checkingShark && !swimmingAnimal.isShark {
	//
	//		}
	//	}
	//}
	bC.lock.Unlock()
}

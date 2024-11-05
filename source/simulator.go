// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-04
// Description:
// This file is for simulating the Wa-Tor Simulation
// According to these rules: https://en.wikipedia.org/wiki/Wa-Tor

package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

// Swimming animal struct used for fish and sharks
type swimmingAnimal struct {
	// Boolean to store whether this is a shark or fish
	isShark bool
	// Counter since last reproduction time
	reproductionTime int
	// Energy counter for starvation
	energy int
	// Used to see if a swimmingAnimal has been checked this iteration
	checkingState bool
}

// Constructor
func newSwimmingAnimal(isShark bool, reproductionTime, energy int) *swimmingAnimal {
	return &swimmingAnimal{isShark: isShark, reproductionTime: reproductionTime, energy: energy}
}

//func (sA *swimmingAnimal) process() {
//	var potentialLocations []int
//
//}

// Primitive data struct to construct bigger chunks
type dataChunk struct {
	data [][]*swimmingAnimal
	// Pointers to surrounding chunks
	leftChunk   *dataChunk
	rightChunk  *dataChunk
	topChunk    *dataChunk
	bottomChunk *dataChunk
}

// Chunk border is used for passing data between threads
type chunkBorder struct {
	afterDataChunk  dataChunk
	beforeDataChunk dataChunk
	lock            sync.Mutex
}

// Constructor
func newChunkBorder(topThreadDataChunk, bottomThreadDataChunk *dataChunk) *chunkBorder {
	afterDataChunk := dataChunk{} // NEED DATA
	beforeDataChunk := dataChunk{}

	afterDataChunk.leftChunk = &afterDataChunk
	afterDataChunk.rightChunk = &afterDataChunk
	afterDataChunk.topChunk = topThreadDataChunk
	afterDataChunk.bottomChunk = &beforeDataChunk

	beforeDataChunk.leftChunk = &beforeDataChunk
	beforeDataChunk.rightChunk = &beforeDataChunk
	beforeDataChunk.topChunk = &afterDataChunk
	beforeDataChunk.bottomChunk = bottomThreadDataChunk

	return &chunkBorder{beforeDataChunk: beforeDataChunk, lock: sync.Mutex{}, afterDataChunk: afterDataChunk}
}

// Process the dataBefore array
func (cB *chunkBorder) processDataBefore(checkingShark bool) {
	cB.lock.Lock()
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
	cB.lock.Unlock()
}

// Handles taking in an integer from the user
func handleInput(inputVariable *int, outputString string) {
	for {
		fmt.Print(outputString)
		_, err := fmt.Scan(inputVariable)
		// Error when not an integer
		if err != nil {
			fmt.Println("Incorrect input: ", err)
			continue
		}
		// Error when inputting values smaller than 0
		if *inputVariable < 0 {
			fmt.Println("Incorrect input: negative value")
			continue
		}
		break
	}
}

// Converts absolute index to X and Y positions
func indexToPosition(index, sizeX int) (posX, posY int) {
	return index % sizeX, index / sizeX
}

// Create chunks of the worldGrid along the X axis
//func createXChunk(xIndex, chunkSizeX int, worldGrid [][]*swimmingAnimal) [][]*swimmingAnimal {
//	var newSlice [][]*swimmingAnimal
//	for index := range worldGrid {
//		newSlice = append(newSlice, worldGrid[index][chunkSizeX*xIndex:(chunkSizeX*xIndex)+chunkSizeX])
//	}
//	return newSlice
//}

func stepTime(chunk [][]*swimmingAnimal, beforeBorder []*swimmingAnimal, afterBorder []*swimmingAnimal) {

}

func main() {
	// Set seed for debugging
	random := rand.New(rand.NewPCG(42, 1024))

	// Declare variables
	//globalCheckingState := false

	var NumShark int
	var NumFish int
	var FishBreed int
	var SharkBreed int
	var Starve int
	var GridSizeX int
	var GridSizeY int
	var Threads int

	// Take in values for variables
	fmt.Println("Type in the following variables:")
	handleInput(&NumShark, "Number of sharks: ")
	handleInput(&NumFish, "Number of fish: ")
	handleInput(&FishBreed, "Number of time units for fish to reproduce: ")
	handleInput(&SharkBreed, "Number of time units for shark to reproduce: ")
	handleInput(&Starve, "Period of time shark can go without food: ")
	handleInput(&GridSizeX, "Size of world, the X dimension: ")
	handleInput(&GridSizeY, "Size of world, the Y dimension: ")
	handleInput(&Threads, "Number of threads to use: ")

	// Cannot use any threads, terminate
	if Threads == 0 {
		fmt.Println("No threads available to be used.")
		return
	}

	totalSize := GridSizeX * GridSizeY
	// Check if there is enough space for all animals
	if NumShark+NumFish > totalSize {
		fmt.Println("Too many animals.")
		return
	}

	// Create grid
	worldGrid := make([][]*swimmingAnimal, GridSizeY)
	for index := range worldGrid {
		worldGrid[index] = make([]*swimmingAnimal, GridSizeX)
	}

	// Create random spawn locations
	positionsSlice := make([]int, totalSize)
	for index := range totalSize {
		positionsSlice[index] = index
	}
	random.Shuffle(totalSize, func(i, j int) {
		positionsSlice[i], positionsSlice[j] = positionsSlice[j], positionsSlice[i]
	})

	// Spawn sharks on grid
	for index := range NumShark {
		newShark := newSwimmingAnimal(true, SharkBreed, Starve)
		posX, posY := indexToPosition(positionsSlice[index], GridSizeX)
		worldGrid[posY][posX] = newShark
	}

	// Spawn fish on grid
	for index := range NumFish {
		newFish := newSwimmingAnimal(false, FishBreed, 0)
		posX, posY := indexToPosition(positionsSlice[index+NumShark], GridSizeX)
		worldGrid[posY][posX] = newFish
	}
	fmt.Println(worldGrid)

	// Split world into chunks for threads
	// Each chunks 2 edges are duplicated
	var chunksSlice [][][]*swimmingAnimal
	swapped := false
	if GridSizeX > GridSizeY {
		GridSizeY, GridSizeX = GridSizeX, GridSizeY
		swapped = true
		//chunkSizeX := GridSizeX / Threads
		//// Get equal sizes
		//for xIndex := range Threads - 1 {
		//	var newSlice [][]*swimmingAnimal
		//	for index := range worldGrid {
		//		newSlice = append(newSlice, worldGrid[index][chunkSizeX*xIndex:(chunkSizeX*xIndex)+chunkSizeX])
		//	}
		//	chunksSlice = append(chunksSlice, newSlice)
		//}
		//var newSlice [][]*swimmingAnimal
		//// Get last one which is a little bigger
		//for index := range worldGrid {
		//	newSlice = append(newSlice, worldGrid[index][chunkSizeX*(Threads-1):])
		//}
		//chunksSlice = append(chunksSlice, newSlice)
	}
	//} else {
	//	chunkSizeY := GridSizeY / Threads
	//	// Get equal sizes
	//	for index := range Threads - 1 {
	//		chunksSlice = append(chunksSlice, worldGrid[chunkSizeY*index:(chunkSizeY*index)+chunkSizeY])
	//	}
	//	// Get last one which is a little bigger
	//	chunksSlice = append(chunksSlice, worldGrid[chunkSizeY*(Threads-1):])
	//}

	chunkSizeY := GridSizeY / Threads
	// Get equal sizes
	for index := range Threads - 1 {
		chunksSlice = append(chunksSlice, worldGrid[chunkSizeY*index:(chunkSizeY*index)+chunkSizeY])
	}
	// Get last one which is a little bigger
	chunksSlice = append(chunksSlice, worldGrid[chunkSizeY*(Threads-1):])
	fmt.Println(chunksSlice)
	ft.println(swapped)

	// Create Atomic chunk border arrays
	if Threads > 1 {
		borderLocks := make([]sync.Mutex, Threads)
		fmt.Println(borderLocks)
		if GridSizeX > GridSizeY {

		} else {

		}
	}

	// Create simulation loop
	for {
		// Start each thread to run its chunk

		// Wait for all threads

		// Stop on no sharks

		// Stop on no fish

		fmt.Println("Run")
		break
	}
}

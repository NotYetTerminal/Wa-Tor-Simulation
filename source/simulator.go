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
	"sync/atomic"
)

// Global variables
var FishBreed int
var SharkBreed int
var Starve int
var GridSizeX int
var GridSizeY int
var Threads int
var CurrentCheckingState = false

// Atomic shared variables for threads
var sharkCount atomic.Int32
var fishCount atomic.Int32
var firstThreadChunk *threadChunk

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

// Sets up the world map for assigning parts to threads
func setUpStartingMap(NumShark, NumFish int) *[][]*swimmingAnimal {
	// Set seed for debugging
	random := rand.New(rand.NewPCG(42, 1024))
	totalSize := GridSizeX * GridSizeY

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
		newShark := newSwimmingAnimal(true, CurrentCheckingState, SharkBreed, Starve)
		posX, posY := indexToPosition(positionsSlice[index])
		worldGrid[posY][posX] = newShark
	}

	// Spawn fish on grid
	for index := range NumFish {
		newFish := newSwimmingAnimal(false, CurrentCheckingState, FishBreed, 0)
		posX, posY := indexToPosition(positionsSlice[index+NumShark])
		worldGrid[posY][posX] = newFish
	}
	fmt.Println(worldGrid[0])
	return &worldGrid
}

// Function ran on all threads simultaneously
func doSimulation(tC *threadChunk, threadCount *atomic.Int32, sharkChan, fishChan chan bool, lock *sync.Mutex, finishedWG *sync.WaitGroup) {
	for {
		// First run for sharks
		tC.processAllRows(true)

		// Wait for all threads
		lock.Lock()
		threadCount.Add(1)
		if threadCount.Load() == int32(Threads) {
			lock.Unlock()
			fmt.Println("Shark: ", sharkCount.Load())
			updateScreen(firstThreadChunk)
			if CurrentCheckingState {
				return
			}
			for range Threads - 1 {
				sharkChan <- true
			}
		} else {
			lock.Unlock()
			<-sharkChan
		}

		// Stop on no sharks
		if sharkCount.Load() == 0 {
			finishedWG.Done()
			return
		}

		// Second run for fish
		tC.processAllRows(false)

		// Wait for all threads
		lock.Lock()
		threadCount.Add(-1)
		if threadCount.Load() == 0 {
			lock.Unlock()
			fmt.Println("Fish: ", fishCount.Load())
			updateScreen(firstThreadChunk)
			// Change checking state before starting a new iteration
			CurrentCheckingState = !CurrentCheckingState
			for range Threads - 1 {
				fishChan <- true
			}
		} else {
			lock.Unlock()
			<-fishChan
		}

		// Stop on no fish
		if fishCount.Load() == 0 {
			finishedWG.Done()
			return
		}
	}
}

func main() {
	// Declare variables
	var NumShark int
	var NumFish int

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

	// Check if there is enough space for all animals
	if NumShark+NumFish > GridSizeX*GridSizeY {
		fmt.Println("Too many animals.")
		return
	}

	worldGrid := setUpStartingMap(NumShark, NumFish)

	// If the X size is bigger than the Y swap them to simplify chunking
	// Save whether swapped to properly render later
	swapped := false
	if GridSizeX > GridSizeY {
		GridSizeY, GridSizeX = GridSizeX, GridSizeY
		swapped = true
	}
	fmt.Println(swapped)

	// Split world into chunks for threads
	var threadChunksSlice []*threadChunk
	chunkSizeY := GridSizeY / Threads
	// Get equal sized chunks
	// They will be the (total rows / threads) - 2 sized because of the border chunks
	for index := range Threads - 1 {
		threadChunksSlice = append(threadChunksSlice, &threadChunk{data: (*worldGrid)[(chunkSizeY*index)+1 : (chunkSizeY*index)+chunkSizeY-1]})
	}
	// Get last one which is a little bigger
	threadChunksSlice = append(threadChunksSlice, &threadChunk{data: (*worldGrid)[chunkSizeY*(Threads-1):]})
	fmt.Println(threadChunksSlice)

	// Create border chunks
	// And connect chunks together
	var borderChunksSlice []*borderChunk
	for index := range Threads {
		tempBorderChunk := newBorderChunk((*worldGrid)[mod((chunkSizeY*index)-1, GridSizeY)], // First row before current threadChunk
			(*worldGrid)[chunkSizeY*index],           // Last row after previous threadChunk
			threadChunksSlice[mod(index-1, Threads)], // Previous threadChunk
			threadChunksSlice[index],                 // Current threadChunk
		)

		borderChunksSlice = append(borderChunksSlice, tempBorderChunk)
		// Update reference in thread chunks
		threadChunksSlice[mod(index-1, Threads)].belowBorderChunk = tempBorderChunk
		threadChunksSlice[index].aboveBorderChunk = tempBorderChunk
	}
	fmt.Println(borderChunksSlice)
	fmt.Println(borderChunksSlice[0].data)

	// WaitGroup for running until finished
	finishedWG := sync.WaitGroup{}
	finishedWG.Add(Threads)

	// Channels for thread synchronisation
	sharkChan := make(chan bool, Threads)
	fishChan := make(chan bool, Threads)

	// Set up atomic variables
	sharkCount.Add(int32(NumShark))
	fishCount.Add(int32(NumFish))
	var threadCount atomic.Int32

	// Mutex for threadCount
	lock := sync.Mutex{}

	// Set first chunk for rendering
	firstThreadChunk = threadChunksSlice[0]

	// Create threads
	for index := range Threads {
		go doSimulation(threadChunksSlice[index], &threadCount, sharkChan, fishChan, &lock, &finishedWG)
	}

	finishedWG.Wait()
	fmt.Println("Simulation done.")
}

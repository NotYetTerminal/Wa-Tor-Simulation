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
func setUpStartingMap(NumShark, NumFish, FishBreed, SharkBreed, Starve, GridSizeX, GridSizeY int) *[][]*swimmingAnimal {
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
	fmt.Println(worldGrid[0])
	return &worldGrid
}

// Function ran on all threads simultaneously
func doSimulation(tC *threadChunk, threadsWGLock *sync.Mutex, threadsWG, finishedWG *sync.WaitGroup) {
	// First run for sharks
	for {
		// Start each thread to run its chunk

		// Each thread runs its border chunks

		// Wait for all threads

		// Stop on no sharks
		break
	}

	// Second run for fish
	for {
		// Start each thread to run its chunk

		// Each thread runs its border chunks

		// Wait for all threads

		// Stop on no fish
		break
	}
}

func main() {
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

	// Check if there is enough space for all animals
	if NumShark+NumFish > GridSizeX*GridSizeY {
		fmt.Println("Too many animals.")
		return
	}

	worldGrid := setUpStartingMap(NumShark, NumFish, FishBreed, SharkBreed, Starve, GridSizeX, GridSizeY)

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
		tempBorderChunk := newBorderChunk(
			(*worldGrid)[mod((chunkSizeY*index)-1, len(*worldGrid))], // First row before current threadChunk
			(*worldGrid)[chunkSizeY*index],                           // Last row after previous threadChunk
			threadChunksSlice[mod(index-1, len(threadChunksSlice))],  // Previous threadChunk
			threadChunksSlice[index],                                 // Current threadChunk
		)

		borderChunksSlice = append(borderChunksSlice, tempBorderChunk)
		// Update reference in thread chunks
		threadChunksSlice[mod(index-1, len(threadChunksSlice))].belowBorderChunk = tempBorderChunk
		threadChunksSlice[index].aboveBorderChunk = tempBorderChunk
	}
	fmt.Println(borderChunksSlice)
	fmt.Println(borderChunksSlice[0].data)

	// WaitGroup for running until finished
	finishedWG := sync.WaitGroup{}
	finishedWG.Add(Threads)

	// WaitGroup for thread synchronisation
	threadsWG := sync.WaitGroup{}
	threadsWG.Add(Threads)
	threadsWGLock := sync.Mutex{}

	// Create threads
	for index := range Threads {
		go doSimulation(threadChunksSlice[index], &threadsWGLock, &threadsWG, &finishedWG)
	}

	finishedWG.Wait()
	fmt.Println("Simulation done.")
}

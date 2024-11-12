// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-04
// Description:
// This file is for simulating the Wa-Tor Simulation
// According to these rules: https://en.wikipedia.org/wiki/Wa-Tor

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"image"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

// Global and Atomic shared variables for threads
var sharkCount atomic.Int32
var fishCount atomic.Int32
var firstThreadChunk *threadChunk
var iteration = 0
var rgbaImage *image.RGBA
var canvasToWrite fyne.CanvasObject
var graphical = true
var swapped = false
var maxIteration = -1
var endTime int64

var FishBreed int
var SharkBreed int
var Starve int
var GridSizeX int
var GridSizeY int
var Threads int
var CurrentCheckingState = false

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
	rand.Shuffle(totalSize, func(i, j int) {
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
			if graphical {
				fmt.Println("Shark:", sharkCount.Load())
				updateScreen(firstThreadChunk)
			}
			for range Threads - 1 {
				sharkChan <- true
			}
		} else {
			lock.Unlock()
			<-sharkChan
		}

		// Second run for fish
		tC.processAllRows(false)

		// Wait for all threads
		lock.Lock()
		threadCount.Add(-1)
		if threadCount.Load() == 0 {
			lock.Unlock()
			if graphical {
				fmt.Println("Fish:", fishCount.Load())
				fmt.Println("Iteration:", iteration)
				updateScreen(firstThreadChunk)
			}
			iteration += 1

			// Change checking state before starting a new iteration
			CurrentCheckingState = !CurrentCheckingState

			// Check for simulation end
			if sharkCount.Load() == 0 || fishCount.Load() == 0 || iteration == maxIteration+1 {
				endTime = time.Now().UnixMilli()
				if iteration != maxIteration+1 {
					fmt.Println("Final iteration:", iteration)
				}
				// Signal simulation ended
				iteration = -1
			} else {
				// Zero counts
				sharkCount.And(0)
				fishCount.And(0)
			}

			for range Threads - 1 {
				fishChan <- true
			}
		} else {
			lock.Unlock()
			<-fishChan
		}

		// Return when simulation ends
		if iteration == -1 {
			finishedWG.Done()
			return
		}
	}
}

func main() {
	fmt.Println("Welcome to Wa-Tor simulation.")

	// Declare variables
	var NumShark int
	var NumFish int

	// Take in values for variables
	developerMode := true
	if developerMode {
		// Preset variables
		NumShark = 200
		NumFish = 400
		FishBreed = 10
		SharkBreed = 60
		Starve = 40
		GridSizeX = 100
		GridSizeY = 100
		Threads = 2

		// Controls whether to show screen and statistical info
		graphical = false

		// Stops the simulation at the set iteration
		maxIteration = 1000
	} else {
		fmt.Println("Type in the following variables:")
		handleInput(&NumShark, "Number of sharks: ")
		handleInput(&NumFish, "Number of fish: ")
		handleInput(&FishBreed, "Number of time units for fish to reproduce: ")
		handleInput(&SharkBreed, "Number of time units for shark to reproduce: ")
		handleInput(&Starve, "Period of time shark can go without food: ")
		handleInput(&GridSizeX, "Size of world, the X dimension: ")
		handleInput(&GridSizeY, "Size of world, the Y dimension: ")
		handleInput(&Threads, "Number of threads to use: ")
	}

	// Start timer
	startTime := time.Now().UnixMilli()

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

	// If the X size is bigger than the Y swap them to simplify chunking
	// Save whether swapped to properly render later
	if GridSizeX > GridSizeY {
		GridSizeY, GridSizeX = GridSizeX, GridSizeY
		swapped = true
	}

	// Must have at least 4 rows
	if GridSizeY < 4 {
		fmt.Println("Grid size must be greater than 3.")
		return
	}

	worldGrid := setUpStartingMap(NumShark, NumFish)

	// Split world into chunks for threads
	var threadChunksSlice []*threadChunk
	chunkSizeY := GridSizeY / Threads
	// Decrease thread count to create correct chunk size
	if chunkSizeY < 4 {
		Threads = GridSizeY / 4
		chunkSizeY = GridSizeY / Threads
	}

	// Get equal sized chunks
	// They will be the (total rows / threads) - 2 sized because of the border chunks
	for index := range Threads - 1 {
		threadChunksSlice = append(threadChunksSlice, &threadChunk{data: (*worldGrid)[(chunkSizeY*index)+1 : (chunkSizeY*index)+chunkSizeY-1]})
	}
	// Get last one which is a little bigger
	threadChunksSlice = append(threadChunksSlice, &threadChunk{data: (*worldGrid)[chunkSizeY*(Threads-1)+1 : len(*worldGrid)-1]})

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

	// WaitGroup for running until finished
	finishedWG := sync.WaitGroup{}
	finishedWG.Add(Threads)

	// Channels for thread synchronisation
	sharkChan := make(chan bool, Threads)
	fishChan := make(chan bool, Threads)

	// Set up atomic variables
	var threadCount atomic.Int32

	// Mutex for threadCount
	lock := sync.Mutex{}

	// Set first chunk for rendering
	firstThreadChunk = threadChunksSlice[0]

	// Set up graphical window
	mainApp := app.New()
	imageWindow := mainApp.NewWindow("Wa-Tor")
	if graphical {
		topLeft := image.Point{}
		var bottomRight image.Point
		if swapped {
			bottomRight = image.Point{X: GridSizeY, Y: GridSizeX}
		} else {
			bottomRight = image.Point{X: GridSizeX, Y: GridSizeY}
		}

		rgbaImage = image.NewRGBA(image.Rectangle{Min: topLeft, Max: bottomRight})
		canvasToWrite = canvas.NewRasterFromImage(rgbaImage)
		imageWindow.SetContent(canvasToWrite)
		if swapped {
			imageWindow.Resize(fyne.NewSize(float32(GridSizeY), float32(GridSizeX)))
		} else {
			imageWindow.Resize(fyne.NewSize(float32(GridSizeX), float32(GridSizeY)))
		}
	}

	// Create threads
	for index := range Threads {
		go doSimulation(threadChunksSlice[index], &threadCount, sharkChan, fishChan, &lock, &finishedWG)
	}

	if graphical {
		imageWindow.ShowAndRun()
	}

	finishedWG.Wait()
	fmt.Println("Simulation done.")
	if sharkCount.Load() == 0 {
		fmt.Println("All sharks dead.")
	} else if fishCount.Load() == 0 {
		fmt.Println("All fish dead.")
	}
	fmt.Println("Total time in seconds:", float32(endTime-startTime)/1000)
}

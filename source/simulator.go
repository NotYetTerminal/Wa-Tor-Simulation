// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-04
// Description:
// This file is for simulating the Wa-Tor Simulation
// According to these rules: https://en.wikipedia.org/wiki/Wa-Tor

package main

import "fmt"

// Swimming animal struct
type swimmingAnimal struct {
	// Boolean to store whether this is a fish or shark
	isFish bool
	// Counter since last reproduction time
	reproductionTime int
	// Energy counter for starvation
	energy int
}

func newSwimmingAnimal(isFish bool, reproductionTime int, energy int) *swimmingAnimal {
	return &swimmingAnimal{isFish: isFish, reproductionTime: reproductionTime, energy: energy}
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

func main() {
	// Declare variables
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

	//fmt.Println(NumShark)
	//fmt.Println(NumFish)
	//fmt.Println(FishBreed)
	//fmt.Println(SharkBreed)
	//fmt.Println(Starve)
	//fmt.Println(GridSizeX)
	//fmt.Println(GridSizeY)
	//fmt.Println(Threads)

	// Create grid
	worldGrid := make([][]swimmingAnimal, GridSizeX, GridSizeY)
	fmt.Println(worldGrid)

	// Spawn fish on grid

	// Spawn sharks on grid

	// Split world into chunks for threads

	// Create Atomic chunk border arrays

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

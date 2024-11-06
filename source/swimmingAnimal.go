// Author: Gábor Major (c00271548@setu.ie)
// Date creation: 2024-11-05
// Description:
// File for storing code related to the swimming animal, which can be shark or fish

package main

// Swimming animal struct used for fish and sharks
type swimmingAnimal struct {
	// Boolean to store whether this is a shark or fish
	isShark bool
	// Used to see if a swimmingAnimal has been checked this iteration
	checkingState bool
	// Counter since last reproduction time
	reproductionTime int
	// Energy counter for starvation
	energy int
}

// Constructor
func newSwimmingAnimal(isShark, checkingState bool, reproductionTime, energy int) *swimmingAnimal {
	return &swimmingAnimal{isShark: isShark, checkingState: checkingState, reproductionTime: reproductionTime, energy: energy}
}

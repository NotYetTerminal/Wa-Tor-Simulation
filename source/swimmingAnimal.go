package main

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

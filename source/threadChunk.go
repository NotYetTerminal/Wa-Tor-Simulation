package main

// Thread chunk is used for holding data for each thread
// Pointers are to surrounding chunks
type threadChunk struct {
	data             [][]*swimmingAnimal
	aboveBorderChunk *borderChunk
	belowBorderChunk *borderChunk
}

// Getters for swimmingAnimal
// May be gotten from a different chunk
func (tC threadChunk) getRightAnimal(indexX, indexY int) *swimmingAnimal {
	return tC.data[indexY][getRightIndex(indexX, len(tC.data[0]))]
}

func (tC threadChunk) getLeftAnimal(indexX, indexY int) *swimmingAnimal {
	return tC.data[indexY][getLeftIndex(indexX, len(tC.data[0]))]
}

func (tC threadChunk) getAboveAnimal(indexX, indexY int) *swimmingAnimal {
	if indexY-1 < 0 {
		return tC.aboveBorderChunk.getBottomRowAnimal(indexX)
	} else {
		return tC.data[indexY-1][indexX]
	}
}

func (tC threadChunk) getBellowAnimal(indexX, indexY int) *swimmingAnimal {
	if indexY+1 > len(tC.data) {
		return tC.belowBorderChunk.getTopRowAnimal(indexX)
	} else {
		return tC.data[indexY+1][indexX]
	}
}

// Setters for swimmingAnimal
// May be set in a different chunk
func (tC threadChunk) setRightAnimal(indexX, indexY int, animal *swimmingAnimal) {
	tC.data[indexY][getRightIndex(indexX, len(tC.data[0]))] = animal
}

func (tC threadChunk) setLeftAnimal(indexX, indexY int, animal *swimmingAnimal) {
	tC.data[indexY][getLeftIndex(indexX, len(tC.data[0]))] = animal
}

func (tC threadChunk) setAboveAnimal(indexX, indexY int, animal *swimmingAnimal) {
	if indexY-1 < 0 {
		tC.aboveBorderChunk.setBottomRowAnimal(indexX, animal)
	} else {
		tC.data[indexY-1][indexX] = animal
	}
}

func (tC threadChunk) setBelowAnimal(indexX, indexY int, animal *swimmingAnimal) {
	if indexY+1 > len(tC.data) {
		tC.belowBorderChunk.setTopRowAnimal(indexX, animal)
	} else {
		tC.data[indexY+1][indexX] = animal
	}
}

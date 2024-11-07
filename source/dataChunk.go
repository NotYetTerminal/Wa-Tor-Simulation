package main

type dataChunk interface {
	getRow(indexY int) []*swimmingAnimal

	getAnimal(indexX, indexY int) *swimmingAnimal
	getLeftAnimal(indexX, indexY int) *swimmingAnimal
	getRightAnimal(indexX, indexY int) *swimmingAnimal
	getAboveAnimal(indexX, indexY int) *swimmingAnimal
	getBelowAnimal(indexX, indexY int) *swimmingAnimal

	setAnimal(indexX, indexY int, animal *swimmingAnimal)
	setLeftAnimal(indexX, indexY int, animal *swimmingAnimal)
	setRightAnimal(indexX, indexY int, animal *swimmingAnimal)
	setAboveAnimal(indexX, indexY int, animal *swimmingAnimal)
	setBelowAnimal(indexX, indexY int, animal *swimmingAnimal)

	getTopRowAnimal(indexX int) *swimmingAnimal
	getBottomRowAnimal(indexX int) *swimmingAnimal

	setTopRowAnimal(indexX int, animal *swimmingAnimal)
	setBottomRowAnimal(indexX int, animal *swimmingAnimal)
}

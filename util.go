package main

func makePositive(number int) int {
	if number < 0 {
		return -number
	}
	return number
}

func makePositiveFloat(number float32) float32 {
	if number < 0 {
		return -number
	}
	return number
}

package main

func makePositive(number int) int {
	if number < 0 {
		return -number
	}
	return number
}

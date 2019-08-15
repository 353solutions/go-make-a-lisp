package main

import (
	"fmt"
)

const π = 3.14159265358

func circleArea(r float64) float64 {
	return π * r * r
}

func main() {
	fmt.Println(circleArea(10))
}

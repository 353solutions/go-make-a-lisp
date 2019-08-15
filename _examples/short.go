package main

import (
	"fmt"
)

func div(a, b int) int {
	return a / b
}

func or(vals ...bool) bool {
	for _, val := range vals {
		if val {
			return true
		}
	}
	return false
}

func main() {
	x := 10

	if x > 4 || div(x, 0) == 0 {
		fmt.Println("OK")
	}

	if or(x > 4, div(x, 0) == 0) {
		fmt.Println("OK")
	}
}

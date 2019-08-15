package main

import (
	"fmt"
)

func collatz(n int) int {
	if n%2 == 0 {
		return n / 2
	}

	return n*3 + 1
}

func main() {
	fmt.Println(collatz(7)) // 22
}

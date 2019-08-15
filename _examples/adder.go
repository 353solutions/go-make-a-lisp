package main

import (
	"fmt"
)

func makeAdder(n int) func(int) int {
	return func(val int) int {
		return val + n
	}
}

func main() {
	add3 := makeAdder(3)
	fmt.Println(add3(4)) // 7
}

package main

import (
	"fmt"
)

func fact(n int) int {
	if n < 2 {
		return 1
	}

	return n * fact(n-1)
}

func main() {
	fmt.Println(fact(10)) // 3628800
}

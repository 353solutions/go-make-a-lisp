package main

import (
	"fmt"
)

func makeAccount(balance float64) func(float64) float64 {
	return func(amount float64) float64 {
		balance += amount
		return balance
	}
}

func main() {
	acct := makeAccount(100)
	fmt.Println(acct(10))
	fmt.Println(acct(-30))
	fmt.Println(acct(70))
}

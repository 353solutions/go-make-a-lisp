package main

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func main() {
	code := []byte("x := collatz(n)")
	fs := token.NewFileSet()
	file := fs.AddFile("", fs.Base(), len(code))
	var s scanner.Scanner

	s.Init(file, code, nil, scanner.ScanComments)

	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%q\n", tok, lit)
	}
}

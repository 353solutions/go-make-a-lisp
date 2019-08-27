package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fmt.Println("Welcome to Hubmle Lisp ☺")
	repl()
	fmt.Println("\nkthxbai")
}

// Token in the language
type Token string

// Tokenize splits the t list of tokens
func Tokenize(code string) []Token {
	code = strings.Replace(code, "(", " ( ", -1)
	code = strings.Replace(code, ")", " ) ", -1)
	var tokens []Token
	for _, tok := range strings.Fields(code) {
		tokens = append(tokens, Token(tok))
	}
	return tokens
}

// Expression to evaluate
type Expression interface{}

// ListExpression - (+ n 1)
type ListExpression struct {
	children []Expression
}

// NumberExpression is a number
type NumberExpression float64

// NameExpression is a name: n, if, ...
type NameExpression string

// ReadExpr reads an expression from slice of tokens
func ReadExpr(tokens []Token) (Expression, []Token, error) {
	var err error
	if len(tokens) == 0 {
		return nil, nil, io.EOF
	}

	tok, tokens := tokens[0], tokens[1:]
	if tok == "(" {
		var children []Expression
		for len(tokens) > 0 && tokens[0] != ")" {
			var child Expression
			child, tokens, err = ReadExpr(tokens)
			if err != nil {
				return nil, nil, err
			}
			children = append(children, child)
		}

		if len(tokens) == 0 {
			return nil, nil, fmt.Errorf("unbalanced expression")
		}

		tokens = tokens[1:] // remove closing ')'
		// FIXME: List expression
	}

	switch tok {
	case ")": // TODO: file:line
		return nil, nil, fmt.Errorf("unexpected ')'")
	}

	// FIXME: number or symbol
	return nil, nil, fmt.Errorf("FIXME")
}

// Read, eval, print, loop
func repl() {
	rdr := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("» ")
		text, err := rdr.ReadString('\n')
		if err != nil {
			break
		}

		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}

		tokens := Tokenize(text)
		fmt.Println("tokens →", tokens)
	}
}

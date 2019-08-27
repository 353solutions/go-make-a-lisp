package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
type Expression interface {
	Eval() (Object, error)
}

// ListExpression - (+ n 1)
type ListExpression []Expression

// Eval implements expression interface
func (e ListExpression) Eval() (Object, error) {
	if len(e) == 0 {
		return nil, fmt.Errorf("empty list expression")
	}

	op, ok := e[0].(NameExpression)
	if ok {
		switch string(op) {
		case "+": // (+ 3 2)
			val := Number(0)
			for _, expr := range e[1:] {
				obj, err := expr.Eval()
				if err != nil {
					return nil, err
				}
				n, ok := obj.(Number)
				if !ok {
					return nil, fmt.Errorf("+ argument not a number: %T", obj)
				}
				val += n
			}
			return val, nil
		}
	}

	return nil, fmt.Errorf("oops")
}

// NumberExpression is a number
type NumberExpression float64

// Eval implements expression interface
func (e NumberExpression) Eval() (Object, error) {
	return Number(e), nil
}

// String implements fmt.Stringer
func (e NumberExpression) String() string {
	return fmt.Sprintf("%f", e)
}

// NameExpression is a name: n, if, ...
type NameExpression string

// Eval implements expression interface
func (e NameExpression) Eval() (Object, error) {
	return nil, fmt.Errorf("not impelemented")
}

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
		return ListExpression(children), tokens, nil
	}

	switch tok {
	case ")": // TODO: file:line
		return nil, nil, fmt.Errorf("unexpected ')'")
	}

	// Number or symbol
	lit := string(tok)
	val, err := strconv.ParseFloat(lit, 64)
	if err == nil {
		return NumberExpression(val), tokens, nil
	}
	return NameExpression(lit), tokens, nil
}

// Object in the language
type Object interface{}

// Number object
type Number float64

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
		//fmt.Println("tokens →", tokens)

		expr, _, err := ReadExpr(tokens)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}
		//fmt.Printf("expr → %#v\n", expr)

		obj, err := expr.Eval()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}
		fmt.Println(obj)
	}
	// (* n 3)
	// (+ (* n 3) 1)
}

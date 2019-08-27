package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	environment = map[string]Object{
		"+": &Function{"+", 0, func(args []Number) (Object, error) {
			total := Number(0.0)
			for _, val := range args {
				total += val
			}

			return total, nil
		}},
		"*": &Function{"*", 0, func(args []Number) (Object, error) {
			total := Number(1.0)
			for _, val := range args {
				total *= val
			}

			return total, nil
		}},
		"%": &Function{"%", 2, func(args []Number) (Object, error) {
			if args[1] == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return Number(int(args[0]) % int(args[1])), nil
		}},
	}
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

	rest := e[1:]
	op, ok := e[0].(NameExpression)
	if ok {
		// speical forms
		switch string(op) {
		case "define": // (define n 3)
			return evalDefine(rest)
		case "or":
			return evalOr(rest)
		}

		obj, err := e[0].Eval()
		if err != nil {
			return nil, err
		}

		fn, ok := obj.(*Function)
		if !ok {
			return nil, fmt.Errorf("bad first expression in list")
		}
		var params []Object
		for _, expr := range rest {
			obj, err := expr.Eval()
			if err != nil {
				return nil, err
			}
			params = append(params, obj)
		}
		return fn.Call(params)
	}

	return nil, fmt.Errorf("oops")
}

func evalDefine(args []Expression) (Object, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("malformed define")
	}
	name, ok := args[1].(NameExpression)
	if !ok {
		return nil, fmt.Errorf("malformed define")
	}
	obj, err := args[2].Eval()
	if err != nil {
		return nil, err
	}
	environment[string(name)] = obj
	return obj, nil
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
	obj, ok := environment[string(e)]
	if !ok {
		return nil, fmt.Errorf("unknown name - %s", e)
	}
	return obj, nil
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

// Function is a function
type Function struct {
	name  string
	nargs int
	op    func(args []Number) (Object, error)
}

func (f *Function) String() string {
	return f.name
}

// Call a function
func (f *Function) Call(args []Object) (Object, error) {
	if f.nargs != 0 && f.nargs != len(args) {
		return nil, fmt.Errorf("%s: wrong number of arguments", f.name)
	}

	var vals []Number
	for i, obj := range args {
		val, ok := obj.(Number)
		if !ok {
			return nil, fmt.Errorf("%s: argument %d of type %T (wanted number)", f.name, i, obj)
		}
		vals = append(vals, val)
	}
	return f.op(vals)
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

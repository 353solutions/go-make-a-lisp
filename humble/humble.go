package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	builtin = &Environment{
		bindings: map[string]Object{
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
		},
	}
)

// Environment holding names
type Environment struct {
	bindings map[string]Object
	parent   *Environment
}

// Find environment holding name
func (e *Environment) Find(name string) *Environment {
	if _, ok := e.bindings[name]; ok {
		return e
	}

	if e.parent == nil {
		return nil
	}
	return e.parent.Find(name)
}

// Get value for name
func (e *Environment) Get(name string) Object {
	return e.bindings[name]
}

// Set a value for name
func (e *Environment) Set(name string, value Object) {
	e.bindings[name] = value
}

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
	Eval(env *Environment) (Object, error)
}

// ListExpression - (+ n 1)
type ListExpression []Expression

func (e ListExpression) String() string {
	s := fmt.Sprintf("%v", []Expression(e))
	// Replace surrounding [] with ()
	return "(" + s[1:len(s)-1] + ")"
}

// Eval implements expression interface
func (e ListExpression) Eval(env *Environment) (Object, error) {
	if len(e) == 0 {
		return nil, fmt.Errorf("empty list expression")
	}

	rest := e[1:]
	op, ok := e[0].(NameExpression)
	if ok {
		// speical forms
		switch string(op) {
		case "define": // (define n 3)
			return evalDefine(rest, env)
		case "or": // (or) -> 0.0, (or 0 2 3) -> 2
			return evalOr(rest, env)
		case "and": // (and) -> 1.0, (and 1 0 2) -> 0
			// FIXME
		case "if": // (if (> x 1) 2 3)
			// FIXME:
		case "lambda": // (lambda (n) (+ n 1))
			return evalLambda(rest, env)
		}

		obj, err := e[0].Eval(env)
		if err != nil {
			return nil, err
		}

		fn, ok := obj.(Callable)
		if !ok {
			return nil, fmt.Errorf("bad first expression in list")
		}
		var params []Object
		for _, expr := range rest {
			obj, err := expr.Eval(env)
			if err != nil {
				return nil, err
			}
			params = append(params, obj)
		}
		return fn.Call(params)
	}

	return nil, fmt.Errorf("oops")
}

func (l *Lambda) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "(lambda (")
	// Can't use strings.Join on []Symbol
	for i, sym := range l.params {
		fmt.Fprint(&buf, sym)
		if i < len(l.params)-1 {
			fmt.Fprint(&buf, " ")
		}
	}
	fmt.Fprintf(&buf, ") ")
	fmt.Fprintf(&buf, "%s", l.body)
	fmt.Fprint(&buf, ")")
	return buf.String()
}

// Callable can be used as a function
type Callable interface {
	Call(args []Object) (Object, error)
}

// (lambda (n) (+ n 1)))
func evalLambda(args []Expression, env *Environment) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("malformed lambda")
	}

	le, ok := args[0].(ListExpression)
	if !ok {
		return nil, fmt.Errorf("malformed lambda")
	}

	params := make([]NameExpression, len(le))
	for i, e := range le {
		s, ok := e.(NameExpression)
		if !ok {
			return nil, fmt.Errorf("malformed lambda")
		}
		params[i] = s
	}
	obj := &Lambda{
		env:    env,
		params: params,
		body:   args[1],
	}
	return obj, nil
}

// Call implements Callable
func (l *Lambda) Call(args []Object) (Object, error) {
	if len(args) != len(l.params) {
		return nil, fmt.Errorf("wrong number of arguments (want %d, got %d)", len(l.params), args)
	}

	m := make(map[string]Object)
	for i, name := range l.params {
		m[string(name)] = args[i]
	}

	env := &Environment{m, l.env}
	return l.body.Eval(env)
}

// Lambda object
type Lambda struct {
	env    *Environment
	params []NameExpression
	body   Expression
}

func evalOr(args []Expression, env *Environment) (Object, error) {
	for _, expr := range args {
		obj, err := expr.Eval(env)
		if err != nil {
			return nil, err
		}
		n, ok := obj.(Number)
		if !ok { // Anything that's not a number is true
			return n, nil
		}
		// Only 0 is false
		if n != 0 {
			return n, nil
		}
	}
	return Number(0.0), nil
}

func evalDefine(args []Expression, env *Environment) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("malformed define")
	}
	name, ok := args[0].(NameExpression)
	if !ok {
		return nil, fmt.Errorf("malformed define")
	}
	obj, err := args[1].Eval(env)
	if err != nil {
		return nil, err
	}
	env.Set(string(name), obj)
	return obj, nil
}

// NumberExpression is a number
type NumberExpression float64

// Eval implements expression interface
func (e NumberExpression) Eval(env *Environment) (Object, error) {
	return Number(e), nil
}

// String implements fmt.Stringer
func (e NumberExpression) String() string {
	return fmt.Sprintf("%f", e)
}

// NameExpression is a name: n, if, ...
type NameExpression string

// Eval implements expression interface
func (e NameExpression) Eval(env *Environment) (Object, error) {
	name := string(e)
	env = env.Find(name)
	if env == nil {
		return nil, fmt.Errorf("unknown name - %s", e)
	}
	return env.Get(name), nil
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

		obj, err := expr.Eval(builtin)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}
		fmt.Println(obj)
	}
	// (* n 3)
	// (+ (* n 3) 1)
}

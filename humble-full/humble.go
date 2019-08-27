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
	builtins *Environment
)

func init() {
	m := map[Symbol]Object{
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
		"begin": &begin{},
		"%": &Function{"%", 2, func(args []Number) (Object, error) {
			if args[1] == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return Number(int(args[0]) % int(args[1])), nil
		}},
		"eq?": &Function{"eq?", 2, func(args []Number) (Object, error) {
			if args[0] == args[1] {
				return 1.0, nil
			}
			return 0.0, nil
		}},
		// MT: In scheme these get arbitrary number of arguments
		"<": &Function{"<", 2, func(args []Number) (Object, error) {
			if args[0] < args[1] {
				return 1.0, nil
			}
			return 0.0, nil
		}},
		"-": &Function{"-", 2, func(args []Number) (Object, error) {
			return args[0] - args[1], nil
		}},
		"/": &Function{"/", 2, func(args []Number) (Object, error) {
			if args[1] == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return args[0] / args[1], nil
		}},
	}

	builtins = &Environment{m, nil}
}

// Token in the language
type Token string

// Tokenize splits the t list of tokens
func Tokenize(code string) []Token {
	code = strings.Replace(code, "(", " ( ", -1)
	code = strings.Replace(code, ")", " )", -1)
	var tokens []Token
	for _, tok := range strings.Fields(code) {
		tokens = append(tokens, Token(tok))
	}
	return tokens
}

// Expression to be computed
type Expression interface {
	Eval(env *Environment) (Object, error)
}

// Object in the language
type Object interface{}

// NumberExpr is a number. e.g. 3.14
type NumberExpr float64

func (e NumberExpr) String() string {
	return fmt.Sprintf("%f", e)
}

// Number is a number in the language
type Number float64

// Eval evaluates value
func (e NumberExpr) Eval(env *Environment) (Object, error) {
	return Number(e), nil
}

// Symbol is a symbol. e.g. pi
type Symbol string

// Eval evaluates value
func (e Symbol) Eval(env *Environment) (Object, error) {
	env = env.Find(e)
	if env == nil {
		return nil, fmt.Errorf("unknown name - %q", e)
	}

	return env.Get(e), nil
}

// ListExpr is a list expression. e.g. (* 4 5)
type ListExpr struct {
	children []Expression
}

func (e *ListExpr) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "(")
	for i, c := range e.children {
		fmt.Fprintf(&buf, "%s", c)
		if i < len(e.children)-1 {
			fmt.Fprintf(&buf, " ")
		}
	}
	fmt.Fprintf(&buf, ")")
	return buf.String()
}

// Eval evaluates value
func (e *ListExpr) Eval(env *Environment) (Object, error) {
	if len(e.children) == 0 {
		return nil, fmt.Errorf("empty list expression")
	}

	rest := e.children[1:]

	// Try special forms first
	op, ok := e.children[0].(Symbol)
	if ok {
		switch op {
		case "define": // (define n 27)
			return evalDefine(rest, env)
		case "set!": // (set! n 27)
			return evalSet(rest, env)
		case "if": // (if (< x 0) 0 x)
			return evalIf(rest, env)
		case "or": // (or), (or 0 1)
			return evalOr(rest, env)
		case "and": // (and), (and 0 1)
			return evalAnd(rest, env)
		case "lambda": // (lambda (n) (+ n 1))
			return evalLambda(rest, env)
		}
	}

	obj, err := e.children[0].Eval(env)
	if err != nil {
		return nil, err
	}

	c, ok := obj.(Callable)
	if !ok {
		return nil, fmt.Errorf("%s is not callable", obj)
	}

	var params []Object
	for _, e := range rest {
		obj, err := e.Eval(env)
		if err != nil {
			return nil, err
		}
		params = append(params, obj)
	}

	return c.Call(params)
}

func evalDefine(args []Expression, env *Environment) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments for 'define'")
	}

	s, ok := args[0].(Symbol)
	if !ok {
		return nil, fmt.Errorf("bad name in 'define'")
	}

	val, err := args[1].Eval(env)
	if err != nil {
		return nil, err
	}
	env.Set(s, val)
	return val, nil
}

func evalSet(args []Expression, env *Environment) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments for 'define'")
	}

	s, ok := args[0].(Symbol)
	if !ok {
		return nil, fmt.Errorf("bad name in 'define'")
	}

	env = env.Find(s)
	if env == nil {
		return nil, fmt.Errorf("unknown name - %s", s)
	}

	val, err := args[1].Eval(env)
	if err != nil {
		return nil, err
	}

	env.Set(s, val)
	return val, nil
}

func evalIf(args []Expression, env *Environment) (Object, error) {
	if len(args) != 3 { // TODO: if without else
		return nil, fmt.Errorf("wrong number of arguments for 'define'")
	}

	cond, err := args[0].Eval(env)
	if err != nil {
		return nil, err
	}

	if cond == 1.0 {
		return args[1].Eval(env)
	}
	return args[2].Eval(env)
}

func evalOr(args []Expression, env *Environment) (Object, error) {
	for _, e := range args {
		obj, err := e.Eval(env)
		if err != nil {
			return nil, err
		}

		val, ok := obj.(Number)
		if !ok {
			return nil, fmt.Errorf("or - %v bad type %T", val, val)
		}

		if val != 0.0 {
			return val, nil
		}
	}

	return Number(0.0), nil
}

func evalAnd(args []Expression, env *Environment) (Object, error) {
	for i, arg := range args {
		obj, err := arg.Eval(env)
		if err != nil {
			return nil, err
		}

		val, ok := obj.(Number)
		if !ok {
			return nil, fmt.Errorf("or - %v bad type %T", val, val)
		}

		if val == Number(0.0) {
			return val, nil
		}

		if i == len(args)-1 {
			return val, nil
		}
	}

	return Number(1), nil
}

func evalLambda(args []Expression, env *Environment) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("malformed lambda")
	}

	le, ok := args[0].(*ListExpr)
	if !ok {
		return nil, fmt.Errorf("malformed lambda")
	}

	params := make([]Symbol, len(le.children))
	for i, e := range le.children {
		s, ok := e.(Symbol)
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

// Callable object
type Callable interface {
	Call(args []Object) (Object, error)
}

type begin struct{}

// Begin function. e.g. (begin (* 3 4) (/ 5 7))
func (b *begin) Call(args []Object) (Object, error) {
	if len(args) == 0 {
		return 0.0, nil
	}

	return args[len(args)-1], nil
}

// Function object
type Function struct {
	name  string
	nargs int
	op    func(args []Number) (Object, error)
}

// Call implement Callable
func (f *Function) Call(args []Object) (Object, error) {
	if f.nargs != 0 && len(args) != f.nargs {
		return nil, f.errorf("wrong number of arguments (want %d, got %d)", f.nargs, len(args))
	}

	var vals []Number
	for i, obj := range args {
		val, ok := obj.(Number)
		if !ok {
			return nil, f.errorf("argument %d: got %v of type %T", i, obj, obj)
		}
		vals = append(vals, val)
	}

	val, err := f.op(vals)
	if err != nil {
		return nil, f.errorf("%s", err)
	}

	return val, nil
}

func (f *Function) errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s - %s", f.name, msg)
}

// Lambda is a lambda object. e.g. (lambda (n) (+ n 1))
type Lambda struct {
	env    *Environment
	params []Symbol
	body   Expression
}

// Call implements Callable
func (l *Lambda) Call(args []Object) (Object, error) {
	if len(args) != len(l.params) {
		return nil, fmt.Errorf("wrong number of arguments (want %d, got %d)", len(l.params), args)
	}

	m := make(map[Symbol]Object)
	for i, name := range l.params {
		m[name] = args[i]
	}

	env := &Environment{m, l.env}
	return l.body.Eval(env)
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
		return &ListExpr{children}, tokens, nil
	}

	switch tok {
	case ")": // TODO: file:line
		return nil, nil, fmt.Errorf("unexpected ')'")
	}

	lit := string(tok)
	val, err := strconv.ParseFloat(lit, 64)
	if err == nil {
		return NumberExpr(val), tokens, nil
	}
	return Symbol(lit), tokens, nil // name
}

// Environment holds name → values
type Environment struct {
	bindings map[Symbol]Object
	parent   *Environment
}

// Find finds the environment holding name, return nil if not found
func (e *Environment) Find(name Symbol) *Environment {
	if _, ok := e.bindings[name]; ok {
		return e
	}

	if e.parent == nil {
		return nil
	}
	return e.parent.Find(name)
}

// Get returns bindings for name in environment
func (e *Environment) Get(name Symbol) Object {
	return e.bindings[name]
}

// Set sets bindings for name
func (e *Environment) Set(name Symbol, value Object) {
	e.bindings[name] = value
}

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
		// fmt.Println("tokens →", tokens)

		expr, _, err := ReadExpr(tokens)
		if err != nil {
			printError(err)
			continue
		}
		//fmt.Printf("expr → %s\n", expr)

		out, err := expr.Eval(builtins)
		if err != nil {
			printError(err)
			continue
		}
		fmt.Println(out)
	}
}

func printError(err error) {
	fmt.Printf("\033[31mERROR: %s\033[0m\n", err)
}

// rlwrap go run humble.go
func main() {
	fmt.Println("Welcome to Hubmle lisp (hit CTRL-D to quit)")
	repl()
	fmt.Println("\nkthxbai ☺")
}

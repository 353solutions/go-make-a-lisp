- log
- font
- silence phone
- slides up to `Code`

~~~
    $ rlwrap go run humble.go
    > (define inc (lambda (n) (+ n 1)))
    > (inc 10)
~~~

- Code
    - Who said code need to be written? (scratch, labview) - 2d
    - Who said code need to be in files? (pharo)
    - Go's source code is UTF-8
- Lex
    - Convert bytes to tokens
    - `go run _examples/toks.go`
    ~~~
    code = b'x = collatz(n)`
    print_tokens(code)
    ~~~
    - show it ignore whitespace
    - $ vim /opt/go/src/go/scanner/scanner.go +688
    - humble.go
    - humble.go:tokenize
	- token is just a str (file, lineno ...)
- Parse
    - $ vim /opt/go/src/go/parser/parser.go +2122
    - `$ go run _examples/ast.go`
    - reader
    - scm.go:readSexpr
    - All numbers are float (like JavaScript)
- Eval
    - story on assember
    - expression vs statement
    - start with just func, name & variable
	- need builtin
	- names + case sensitive
	- lisp 1/2
    - run & repl
	- fmt.Stringer
	- readline for matching ()
    - if
	- what are booleans? (we'll use 1.0 & 0.0)
	- short circuit
~~~
func div(a, b int) int {
	return a / b
}

func or(vals ...bool) bool {
	for _, val := range vals {
		if val {
			return true
		}
	}
	return false
}

func main() {
	x := 10

	if x > 4 || div(x, 0) == 0 {
		fmt.Println("OK")
	}

	if or(x > 4, div(x, 0) == 0) {
		fmt.Println("OK")
	}
}
~~~
    - or, and
	- they do it
	- `(or)` -> 0.0
	- `(and)` -> 1.0
	- `True or 1/0`
    - define & set!
    - lambda
	- parameter passing (value, ref ..)
	- square.scm
	- collatz.scm
	- fact.scm (recursion)
	- adder.scm (closure)
	    - Show Go
    - account.scm
    - begin

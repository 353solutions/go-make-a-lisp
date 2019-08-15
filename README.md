# Go Make a Lisp Interpreter

[Workshop](https://www.meetup.com/Go-Israel/events/263973009/) given as part of
[Go Israel](https://www.meetup.com/Go-Israel/)

"If you don't know how compilers work, then you don't know how computers work."
    - Steve Yegge

We'll implement a small [lisp like](http://norvig.com/lispy.html) language and
discuss language design & implementation  issues and how they are found in
Go.

- [Lexing](https://en.wikipedia.org/wiki/Lexical_analysis) & Parsing: What are
  the implication of Go's C based syntax
- Variable scope &
  [closures](https://en.wikipedia.org/wiki/Closure_(computer_programming)):
  Understand some common scope issues in Go
- Types: Why do we have several number types? Other types in Go
- Evaluating code: Why does `or` and `and` [short
  curcuit](https://en.wikipedia.org/wiki/Short-circuit_evaluation)

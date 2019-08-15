(define collatz
  (lambda (n)
    (if (eq? (% n 2) 0)
	(/ n 2)
	(+ (* n 3) 1))))


(collatz 7) ; 22

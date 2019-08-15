(define make-adder
  (lambda (n)
    (lambda (val)
      (+ val n))))

(println ((make-adder 3) 4)) ; 7

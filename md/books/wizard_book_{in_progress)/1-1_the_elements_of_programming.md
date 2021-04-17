# 1. Building abstractions with procedures

## 1.1 The elements of programming

Powerful languages have three mechanisms to combine simple ideas into complex
ones:

1. Primitive expressions
2. Means of combination
3. Means of abstraction

We deal with 'two'[^1] elements: procedures and
data.

[^1]: They're apparently not that distinct.

Data is stuff. Procedures are what we do to the stuff.

This chapter is focused only on the rules for building procedures.

### 1.1.1 Expressions

Lisp uses **prefix notation**, i.e. (+ 1 3) not 1 + 3. We know this from Racket
because we cheated by doing some Racket first.

Use pretty-printing and don't be rude.

The interpreter:

- reads the expression from a terminal
- evaluates it
- prints the results

This is what a read-eval-print loop is! That makes sense now!

### 1.1.2 Naming and the environment

We can name stuff with `define`, i.e.

```lisp
(define size 2)
```

The interpreter then knows that when we say `size` what we mean is 2.

The book says `define` is the simplest means of abstraction --- which makes
sense. For example

```lisp
(define pi 3.14159)
(define radius 10)
(define circumference (* 2 pi radius))
```

That's quite a lot of abstraction!

> Indeed, complex programs are constructed by building, step by step,
> computational objects of increasing complexity.

If we're associating values with symbols, that must be maintained in some
memory. This is called the _global environment_.

### 1.1.3 Evaluating combinations

A goal of this chapter is "to isolate issues about thinking procedurally" ---
let's see if we can do that too!

When evaluating combinations, the interpreter will:

1. Evaluate all the subexpressions
2. Apply the procedure that the leftmost subexpression (the operator) to all
   the other subexpression values (the operands)

This is a **recursive process**: the first step is to do the process again, on
a deeper level.

Introduces this _tree_ notation where values "percolate upwards" --- this is an
example of **tree accumulation**, although I don't know what that is.

The book specifies that for example, `define` is **not a combination**. We
don't "apply `define` to `x` and 3" or whatever. This is an example of a
**special form**: it has its own evaluation rules.

### 1.1.4 Compound procedures

**Procedure definitions** are an even more powerful abstraction technique. What
this means is, giving a name to some compound operation and then being able to
refer to it by that name.

For example:

```lisp
(define (square x) (* x x))
```

We've abstracted an expression (multiplying a number by itself) into the name
`square`.

The general form is

```lisp
(define (<name> <formal parameters>) <body>)
```

Why is this powerful?

```lisp
(define (sum-of-squares x y) (+ (square x) (square y)))
```

We've abstracted --- from our abstraction! And so on, and so forth. These
abstractions are the building blocks of our lovely house.

### 1.1.5 The substituion model for procedure application

The interpreter has a process for evaluating a combination with an abstracted,
compound procedure too. That is:

1. Evaluate the elements of the combination
2. Apply the procedure to the arguments

The book goes into how to do this using an example. The specific method is
called the **substitution model**: find the named procedure, substitute it,
continue.

Ooh fun we're going to implement an interpreter and compiler!

They emphasise that this is a **model**, an incomplete picture of what's
actually going on. There's more detailed cases where this isn't quite right.

Applicative order is when you evaluate combinations as they appear, rather than
expanding everything out and then reducing (which may seem the more natural way
to do it). In a lot of cases they're the same, but they aren't _always_ the
same.

### 1.1.6 Conditional expressions and predicates

If we want to output certain values based on an input, this is called a _case
analysis_. For example, computing the absolute value of a number requires some
case analysis.

To do this, we need conditional logic.

The general form (I'm going to use square brackets because I'm not rude) is

```lisp
(cond [<p1> <e1>]
	  [<p2> <e2>]
	  ...
	  [<pn> <en>])
```

These `[<p> <e>]` pairs are called clausesL we have a predicate (if this) and
an expression (then this).

The interpreter evaluates these in order: it checks if `p1` is true, then `p2`,
and so on.

Formally, a **predicate** is a procedure that returns either true or false, as
well as any expressions that evaluate to true or false. Think "predicated
upon" = "conditional upon".

There is also an `else` clause, which just means if nothing else evaluates to
true then do this.

If there are precisely two cases in the case analysis, we can also use an
if-conditional:

```lisp
(define (abs x)
	(if (< x 0)
		(- x)
		x))
```

There are also logical composition operators, most commonly `and`, `or` and
`not`.

## Exercises

Going to take a quick break from reading to do some maths!

### Exercise 1.1

```lisp
10 => 10
(+ 5 3 4) => 12
(- 9 1) => 8
(/ 6 2) => 3
(+ (* 2 4) (- 4 6)) => (+ 8 -2) => 6
(define a 3) => 3
(define b (+ a 1)) => 4
(+ a b (* a b)) => (+ 3 4 (* 3 4)) => 19
(= a b) => #f

; serious mode
(if (and (> b a) (< b (* a b)))
	b
	a)
=>
(if (and (> 4 3) (< 4 12)))
	4
	3)
(if (and #t #t)
	4
	3) => 4

(cond ((= a 4) 6)
	  ((= b 4) (+ 6 7 a))
	  (else 25))
=> 16

(+ 2 (if (> b a ) b a))
=>
(+ 2 (if #t 4 3)) => 6

(* (cond ((> a b) a)
		 ((< a b) b) ; this one is true
		 (else -1))
	(+ a 1))
=>
(* 4 4) => 16
```

These were all right I am basically the computer.

### Exercise 1.2

Translating an expression into prefix form.

```lisp
(/ (+ 5 4 (- 2 (- 3 (/ 6 (/ 4 5))))) (* 3 (- 6 2) (- 2 7)))
```

### Exercise 1.3

Take three numbers as arguments, and return the sum of the squares of the two
larger numbers

```lisp
(define (sum-of-squares-largest a b c)
	(cond [(and (< a b) (< a c)) (+ (* b b) (* c c))]
		  [(and (< b a) (< b c)) (+ (* a a) (* c c))]
		  [else 				 (+ (* a a) (* b b))]))
```

### Exercise 1.4

This is kind of spicy. Describe the behaviour of this procedure

```lisp
(define (a-plus-abs-b a b)
	((if (> b 0) + -) a b))
```

If `b` is positive, it returns `a + b`. If `b` is _negative_, it returns `a -
b`. The if-statement returns an **operator**, not a **value**. This is very
cool and I didn't know this worked.

So in effect, it adds the absolute value of `b` to `a` (since |b| = -b if b <
0).

I guess this is an example of what they meant when they said 'data' and
'procedures' are kind of the same thing: an operation can be conditional, and
can be manipulated like data.

### Exercise 1.5

This kind of hurts my brain so this is where I'm going to finish for tonight.

Let's think like Lisp:

```lisp
(define (p) (p))

(define (test x y)
	(if (= x 0)
		0
		y))


(test 0 (p))
=>
(if (= 0 0)
	0
	(p))
=> 0 			; first condition is true
```

This is applicative order. We never read (p), so we never have to evaluate it
(this is good because (p)'s definition is weird and recursive.

Let's now think like a dumb human being. We expand everything, and then we
determine the output.

OK, let's go:

```dumb-meat
(test 0 (p))
=>
(if (= 0 0)
	0
	(p))
=> expand p which is a named procedure of some kind
(if (= 0 0)
	0
	(p))
=> hang on this is the same oh no oh no
(if (= 0 0)
	0
	(p))
```

and so on. The fact that the predicate is true, and that we don't _have_ to
evaluate the alternative expression, doesn't matter. We expand everything, so
we get stuck in a loop.

OK this isn't quite right. I think my _computer_ explanation was wrong. When
an applicative-order evaluator sees `(p)`, it starts to try to expand it
immediately. 

I don't really understand.

---

So I had a read --- I think I basically got things the wrong way around! So I
_am_ cleverer than a machine.

When an interpreter using applicative-order evaluation sees `(test 0 (p))`, it:

1. Evaluates the subexpressions
2. Applies the operation

So it goes:

1. OK applying test to 0 and (p), what do 0 and (p) mean?
2. I know 0, 0=0 that's easy
3. OK what does (p) mean? Ah! (p) = (p)! If only I knew what _(p)_ was.
4. Wait, there's a definition! (p) = (p)! Now all that's left to do is evaluate
   (p)...

&c.

Whereas I, in my infinite wisdom, can say, "Well it doesn't matter what (p)
evaluates to because we're never going to get there are we?! (= 0 0) and that's
the end of that!"

And on that note I think I can sleep.

**Note:** Hannah, start back from
[here](https://mitpress.mit.edu/sites/default/files/sicp/full-text/book/book-Z-H-10.html#%25_sec_1.1.7).

---

And we're back.

### 1.1.7 Example: square roots by Newton's method

We're talking about the difference between declarative knowledge and imperative
knowledge. For example, I can say that the square root of 5 is a number that,
when squared, gives you 5. That doesn't tell you how to calculate it.
Equivalently, I _know_ there are infinitely many prime numbers without having
to tell you how to calculate them in some way.

There is a difference between describing things, and describing how to _do_
things.

The book gives the example of calculating a square root using Newton's
approximation.

### Exercise 1.6

I think this is a problem of applicative-order evaluation. Since `new-if` isn't
a special form, the interpreter will try to:

1. Evaluate the sub-expressions
2. Apply the operation

Since the else-clause recursively calls out function, the interpreter will get
in a loop trying to evaluate sub-sub-expressions, and so on.

### Exercise 1.7

Let the number we want to calculate the square root of be such that its square
root is less than 0.001. Say, 0.0001^2. Then since our initial guess was
greater than 0.001, the iterative process will terminate at just below 0.001.
This is because the actual solution is less than the tolerance.

Example: 

```lisp
(square (sqrt 0.00001)) ~= 0.001
```

So we're out by an order of 100.

Equivalently with very large numbers. Let the number be large enough that the
difference between it and its approximation is greater than 0.001. Then the
process cannot terminate, since `good-enough?` is never applied.

Example:

```lisp
(sqrt (expt 10 64))
```

will never terminate (on my computer).

Alternative:

```lisp
(define (good-enough previous-guess guess)
	(< (abs (- previous-guess guess)) (/ previous-guess 100)))
```

This calculates both the above to some precision.

### Exercise 1.8

The only difficulty here is me getting my x's and y's confused.

### 1.1.8 Procedures as black-box abstractions

The `sqrt` program can be considered as a cluster of procedures, some dependent
on each other. This tree of dependencies mirrors the decomposition of the
problem (finding a square root) into sub-problems (is my guess good enough,
etc.).

Each procedure **accomplishes a task** that can be used to define other
procedures. The book calls this a **procedural abstraction**: we don't care how
`good-enough?` does its job, particularly, when we're looking at a higher
level.

Procedure definitions suppress detail.

Talking about local names. The specific names of the formal parameters of a
procedure should not affect the user of the procedure.

This is called a **bound variable**. We say that it is **locally scoped** to
its procedure. The alternative to being bound, is being free.

We can also **localise the sub-procedures themselves**, for example if we have
a large codebase and expect someone to use the name `good-enough?` later on.
This is called **block structure**, and this whole discipline of scoping
formal parameters and procedures is referred to as **lexical scoping**.

This allows us to simplify things a lot: a local procedure is "aware" of a
local formal parameter, so we can use it without having to pass it in as one.

EOF

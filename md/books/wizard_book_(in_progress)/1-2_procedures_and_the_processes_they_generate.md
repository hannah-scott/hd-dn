## 1.2 Procedures and the processes they generate

We've covered the basic 'rules' of the game and now we need to know some
strategy, i.e. common patterns of usage.

This chapter is focused on common "shapes" that simple processes follow, as
well as some information on computational resources.

### 1.2.1 Linear recursion and iteration

Looking at the factorial function as an example. If we specify that 1! = 1, and
that n! = n(n - 1)!, this defines our procedure _recursively_.

Alternatively, we can maintain a running product and a counter. While the
counter is less than n, we set the product to product * counter and we
increment the counter.

What is the difference?

If we evaluate the first one, we see an expansion and contraction:

	(fact 3)
	(* 3 (fact 2)
	(* 3 (* 2 (fact 1)
	(* 3 (* 2 1)
	(* 3 2)
	6

This expansion (the number of terms of which is linear to n) is called **linear
recursion**.

If we expand the second, we see:

	(fact-iter 1 1 3)
	(fact-iter 1 2 3)
	(fact-iter 2 3 3)
	(fact-iter 6 4 6)
	6

Here there is no expansion. The terms still grow linearly with n. This is
called **linear iteration**.

The difference is in the underlying shape, as well as what the interpreter has
to keep track of. For a linear recursive process, the amount of information
needed grows linearly with n (since the interpreter needs to keep track of all
the nested operations). So do the number of steps. But for a linear _iterative_
process, the number of "state variables" is fixed and only the number of steps
grow with n.

This is where the terminology is confusing. Both of the above are recursive
**procedures** (procedures that refer to themselves), but when we say a
**process** is recursive we're talking about how the process of evaluating the
procedure evolves over time.

In some languages this is purely academic: the interpreter evaluates _all_
recursive procedures with the memory footprint of a recursive process. This is
true of C, but not of Scheme. An implementation of a recursive procedure that
evaluates as an iterative process in constant space is called
**tail-recursive**.

### Exercise 1.9

First:

	(define (+ a b)
		(if (= a 0)
			b
			(inc (+ (dec a) b))))

	(+ 4 5)
	(inc (+ 3 5))
	(inc (inc (+ 2 5)))
	(inc (inc (inc (+ 1 5))))
	(inc (inc (inc (inc (+ 0 5)))))
	(inc (inc (inc (inc (5)))))
	(inc (inc (inc (6))))
	(inc (inc (7)))
	(inc 8)
	9

This is a linearly recursive process.

Second (I assume this is linearly iterative but let's prove it):

	(define (+ a b)
		(if (= a 0)
			b 
			(+ (dec a) (inc b))))

	(+ 4 5)
	(+ 3 6)
	(+ 2 7)
	(+ 1 8)
	(+ 0 9)
	9

So this process is linearly iterative.

**Note:** Um linearly iterative processes are a lot nicer to evaluate by
hand... be kind to your computer and use them.

### Exercise 1.10

Ackermann's function. Let's get some paper.

	(A 1 10)
	(A 0 (A 1 9))
	(A 0 (A 0 (A 1 8)))
	...
	(A 0 nine times (A 1 1))
	(A 0 eight times (A 1 2))
	(A 0 seven times (A 1 (* 2 2)))
	...
	1024

Each nested `(A 0 y*)` multiplies `y*` by 2. The final term is 2. So we have
2^10. 

	(A 2 4)
	(A 1 (A 2 3)
	(A 1 (A 1 (A 2 2)))
	(A 1 (A 1 (A 1 (A 2 1))))		; y = 1 -> 2
	(A 1 (A 1 (A 1 2)))
	(A 1 (A 1 (A 0 (A 1 1))))
	(A 1 (A 1 (A 0 2)))				; x = 0 -> 4
	(A 1 (A 1 4))
	(A 1 (A 0 (A 1 3)))
	(A 1 (A 0 (A 0 (A 1 2))))
	(A 1 (A 0 (A 0 (A 0 (A 1 1))))) ; y = 1 -> 2
	(A 1 (A 0 (A 0 (A 0 2))))		; x = 0 -> 4
	(A 1 (A 0 (A 0 4)))				; x = 0 -> 8
	(A 1 (A 0 8))					; x = 0 -> 16
	(A 1 16)						; from the above				
	2^16

So

	(A 2 n)

evaluates to

	(A 1 (A 1 ( ... ( A 1 1))))

n times, or

	2 ^ n

Let's look at three, but I'm fairly sure we have a tower of power:

	(A 3 3)
	(A 2 (A 3 2))
	(A 2 (A 2 (A 3 1)))				; y = 1 -> 2
	(A 2 (A 2 2))
	(A 2 2^2)						; from the above
	2 ^ ( 2 ^ ( 2^2))
	2 ^ 2 ^ 2 ^ 2 

This is gross. My poor computer.

**Note:** Having thought about this, 2^16 = 65536. You can calculate this by
because 2^10 = 1024 and 2^6 = 64, or if you prefer because ints are between
-65535 and +65535 on 32-bit architecture (16 binary bits either way).

Let's look at concise mathematical definitions (we love maths here).

	(define (f n) (A 0 n))
	(define (g n) (A 1 n))
	(define (h n) (A 2 n))
	(define (k n) (* 5 n n))

We're given k(n) = 5n^2.

From the definition, f(n) = 2n if n != 0, 0 if n = 0. But 2 * 0 = 0, so f(n) =
2n.

We have g(n) = 2 ^ n. Why?

	(A 1 n)
	(A 0 (A 1 (- n 1)))
	...
	(A 0 ... n-1 times ... (A 1 1))
	(A 0 ... n-2 times (A 0 2))

And each (A 0 k) evaluates to 2k. So we multiply by 2 each time, 2^n.

Similarly, (A 2 n) is 2 ^ 2 ^ ... ^ 2, n times. This is tetration, let's write
^n 2.

### 1.2.2 Tree recursion

Here the example is an implementation of the Fibonacci series:

	(define (fib n)
		(cond ((= n 0) 0)
			  ((= n 1) 1)
			  (else (+ (fib (- n 1))
					   (fib (- n 2))

The process evolves like a tree: to get (fib 5), we need 4 and 3. To get 3, we
need 2 and 1, etc. Lots of branching.

The number of steps to calculate this process grows exponentially, but the
space required is proportional to the maximum depth of the tree.

The Fibonacci numbers could be implemented as a linearlly iterative process
which would run a lot faster and use less space.

However, there are some cases when tree recursion is useful.

There is a challenge to write an iterative version of a change-counting
algorithm.

### Exercise 1.11

As a recursive process:

	(define (f n)
		(cond [(< n 3) n]
			  [else (+ (f (- n 1))
					   (* 2 (f (- n 2)))
					   (* 3 (f (- n 3))))]))

And as a linearly iterative process:
	
	(define (f2 n)
	  (f-iter 2 1 0 n))
	
	(define (f-iter a b c count)
	  (if (= count 0)
	      c
	      (f-iter (+ a (* 2 b) (* 3 c)) a b (- count 1))))

Unsurprisingly the linearly iterative version is a _lot_ quicker.

### Exercise 1.12

Generate the elements of Pascal's triangle with a recursive process.

We're going to take this to mean, given a row and column, calculate that row
and column of Pascal's triangle.

Let's start with some degenerate cases: if r < 1 or c < 1, return 0. If c > r,
return 0 (this lets us add along the right edge). If r = 1 and c = 1, return 1.

In the other case, we want to add the number above and to the left (r - 1, c -
1), and above and to the right (r - 1, c).

So maybe:
	
	(define (pascal r c)
	  (cond [(or (< r 0) (< c 0) (< r c)) 0]
	        [(and (= r 1) (= c 1)) 1]
	        [else (+ (pascal (- r 1) (- c 1)) (pascal (- r 1) c))]))

This works.

### Exercise 1.13

Double proof by induction with a lemma!

### 1.2.3 Orders of growth

We're using big-theta notation to denote orders of growth. By an order of
growth, we mean some important resource (number of operations, storage used) is
growing with some important feature of our parameters (number of digits of
accuracy, size of input, rows in a matrix).

This is obviously a very crude measure since you can be out by any integer
multiple, but it's a good guide for long-term performance.

### Exercise 1.14

I drew it on some paper.

As the amount to be changed grows, the number of nodes or decision points is
going to grow like n^2.

Think of it like an x/y graph. The x-axis is (cc amount (- kinds-of-coins 1))
--- our amount doesn't change. This is fixed at length 5, the number of
different coins we have.

Meanwhile, the y-axis is (cc (- amount (first-denomination kinds-of-soins))
kinds-of-coins). As the amount of money grows, the number of times we have to
reduce it by say 50 cents to get to below 0 increases linearly. That is, n/50
is linear in n (obviously).

All the other nodes are in the middle and none are trimmed, so number of nodes
is going to be n^2. So number of steps is also n^2.

We can think of space as like the "width" of the "bottom" row (without knowing
what that means). If we expand everything, that's the worst case scenario for
space. That is like n, so spaces grows linearly with n.

### Exercise 1.15

(a) Since we are not computers, we can ignore the details of p.

We know p will recursively be called until |angle| < 0.1. At each step, we will
do angle -> angle/3. So effectively, find n such that 12.15/3^n < 0.1.

Equivalently, 12.15/0.1 = 121.5 < 3^n => n = 5.

So p will be applied 5 times in the evaluation of `(sine 12.15)`.

(b) I think this grows logarithmically in both steps and space.

We know from (a) that the process takes n "steps" (calls of p) when a/3^n <
0.1.

So
	a / 3^n < 0.1 <=>
	10 a < 3^n <=>
	k log(a) < n,

where k is some constant. So the number of steps is logarithmic, and the amount
of expansion required is also logarithmic so the space is too.

### 1.2.4 Exponentiation

Shown is a neat logarithmic-time exponentiation algorithm.

It makes use of the fact that for n an even number, x^n = (x ^(n/2))^2.

Of note is the method for showing that it is logarithmic time: if we double the
input, then we require only one extra step to calculate the answer. So
multiplying the input increases the number of steps additively, which is what
we mean by logarithmic really.

Right --- I'm going to take a break here because that was a lot of maths.

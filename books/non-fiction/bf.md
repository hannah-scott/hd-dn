# Beautiful Racket: `bf`

## 1. Intro

`bf` is a minimalist language that uses 8 characters to perform simple
operations on a array of bytes, reading and writing them to `stdin` and
`stdout`.

Every language in Racket has a **reader**, which converts the source to
S-expressions; and an **expander**, which turns those S-expressions into Racket
code.

If our reader can generate "clean, well-structured S-expressions" from source,
that will make the work of writing a clean expander much easier --- and vice
versa.

## 2. Grammars and Parsers

<aside>I always misspell "grammars" as "grammers" so apologies in
advance.</aside>

A **parse tree** is a hierarchical representation of source code that describes
its "structure" in the sense of e.g. recursion, conditional statements: a
function which generates them is called a **parser**.

One can produce a specification of language structure (called a **grammar**) and
from this use a parser generator --- this has advantages over hand-coding the
parser, e.g. grammar changes.

A grammar can also act as a "reality check":

> If we can write a grammar, chances are good that we can make our language
> work. If we can't, it's a warning.

## 3. Grammar Notation

A grammar is a set of **production rules**: a name saying what the structural
element is called in our language; and a pattern for what that element looks
like (similar to a regular expression).

The parser attempts to match our source code to the rules of our grammar,
starting from the first rule and working recursively through them until either:

1. It reaches **terminals**, elements which cannot be decomposed further, and
   succeeds
2. It can't figure out how to reach terminals, and fails

A grammar can be **recursive** and still well-defined: the example given is a
subset of S-expressions, which might contain elements of that subset, e.g.

```lisp 
(+ 1 (* 2 3)) 
```

## 4. The Parser

The structure of `bf` turns out to be quite simple: each program is made up of
either operations, or a loop which contains either operations or loops.

There are alternative grammars for our language, e.g. we could say that a
program contains either operations, or loops which contain programs.

The difference between these two grammars comes down to ease of writing the
eventual expander, and limiting the number of grammar rules.

## 5. The Tokenizer and Reader

A **token** is the "smallest meaningful chunk" of source code, e.g. each
character within the code: a function called a **tokenizer** converts source to
tokens.

A tokenizer can be useful if it reduces the number of distinct things we need to
create rules for, e.g.

1. We may simplify our lives by remove all comments
2. Decimals, fractions, etc. can all be handled as just "numbers"

There are downsides too:

1. Removed strings e.g. comments become completely invisible to our parser,
   which we may not want
2. Tokens are like building blocks: making them too large may limit what we can
   make



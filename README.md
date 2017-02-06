# bnf

This is a small set of Go tools for dealing with context free grammars, and in particular those represented in [Backus-Naur form](https://en.wikipedia.org/wiki/Backus–Naur_form).

For more information on the parser API, see the [GoDoc](https://godoc.org/github.com/unixpickle/bnf).

# BNF example

To get a sense for the BNF syntax, here is a BNF for quoted strings containing the letters 'a', 'b', and 'c':

```bnf
<qstr> ::= '"' <str> '"'
<str>  ::= "" | "a" <str> | "b" <str> | "c" <str>
```

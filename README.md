WIP
---

This project is a *work in progress*. The implementation is *incomplete* and
subject to change. The documentation may be inaccurate.

asm
===

[![Build Status](https://travis-ci.org/USER/REPO.svg?branch=master)](https://travis-ci.org/mewlang/asm)
[![Coverage Status](https://img.shields.io/coveralls/mewlang/asm.svg)](https://coveralls.io/r/mewlang/asm?branch=master)
[![GoDoc](https://godoc.org/github.com/mewlang/asm?status.svg)](https://godoc.org/github.com/mewlang/asm)

The aim of this project is to create a specification of the assembly language
using [EBNF][]. Based on this specification a lexer and a parser will be
developed along with an assembler and a disassembler.

[EBNF]: https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_Form

Documentation
-------------

Documentation provided by GoDoc.

- [lexer][]: implements tokenization of assembly source text.
- [token][]: defines constants representing the lexical tokens of the assembly
language.

[lexer]: http://godoc.org/github.com/mewmew/asm/lexer
[token]: http://godoc.org/github.com/mewmew/asm/token

Examples
--------

### tokens

The [tokens][examples/tokens] command demonstrates how to tokenize input files
using the [Parse][lexer.Parse] function.

	go get github.com/mewmew/asm/examples/tokens

[examples/tokens]: https://github.com/mewmew/asm/blob/master/examples/tokens/tokens.go#L23
[lexer.Parse]: http://godoc.org/github.com/mewmew/asm/lexer#example-Parse

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/

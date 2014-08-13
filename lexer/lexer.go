// The implementation of this package is heavily inspired by Rob Pike's amazing
// talk titled "Lexical Scanning in Go" [1].
//
// The assembly language specification have been based upon the Go Programming
// Language Specification [2] and some comments are therefore governed by a
// BSD-style license [3]. Any original content is hereby released into the
// public domain [4].
//
// [1]: https://www.youtube.com/watch?v=HxaD_trXwRE
// [2]: http://golang.org/ref/spec
// [3]: http://golang.org/LICENSE
// [4]: https://creativecommons.org/publicdomain/zero/1.0/

// Package lexer implements tokenization of assembly source code.
package lexer

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mewmew/asm/token"
)

const (
	// eof is the rune returned by next when no more input is available.
	eof = -1
)

// A lexer lexes the input into tokens.
type lexer struct {
	// The input text.
	input string
	// The start position of the current token.
	start int
	// The current position in the input.
	pos int
	// The width in byte of the last rune read with next.
	width int
	// A slice of scanned tokens.
	tokens []token.Token
}

// Parse lexes the provided input and returns a slice of scanned tokens.
func Parse(input string) (tokens []token.Token) {
	l := &lexer{
		input:  input,
		tokens: make([]token.Token, 0),
	}

	// Tokenize the input.
	l.run()

	return l.tokens
}

// run lexes the input by repeatedly executing the active state function until
// it returns a nil state.
func (l *lexer) run() {
	// lexLine is the initial state function of the lexer.
	for state := lexLine; state != nil; {
		state = state(l)
	}
}

// errorf emits an error token and terminates the scan by returning a nil state
// function.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	tok := token.Token{
		Typ: token.Error,
		Val: fmt.Sprintf(format, args...),
	}
	l.tokens = append(l.tokens, tok)
	return nil
}

// emitEOF emits an EOF token and terminates the scan by returning a nil state
// function.
func (l *lexer) emitEOF() stateFn {
	if l.pos < len(l.input) {
		log.Fatalf("lexer.emitEOF: unexpected eof; pos %d < len(input) %d.\n", l.pos, len(l.input))
	}
	if l.start != l.pos {
		log.Fatalf("lexer.emitEOF: invalid eof; pending input %q not handled.\n", l.input[l.start:])
	}
	l.emit(token.EOF)
	return nil
}

// emit emits a token of the provided token type and advances the token start
// position.
func (l *lexer) emit(typ token.Type) {
	tok := token.Token{
		Typ: typ,
		Val: l.input[l.start:l.pos],
	}
	l.tokens = append(l.tokens, tok)
	l.start = l.pos
}

// next consumes the next rune of the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune of the input.
func (l *lexer) peek() (r rune) {
	r = l.next()
	if r == eof {
		return eof
	}
	l.backup()
	return r
}

// backup backs up one rune in the input. It can only be called once per call to
// next.
func (l *lexer) backup() {
	if l.width == 0 {
		// TODO(u): Handle eof elsewhere so we never hit this case.
		log.Fatalln("lexer.backup: invalid width; no matching call to next.")
	}
	l.pos -= l.width
	l.width = 0
}

// isLetter returns true if r is a Unicode letter or an underscore, and false
// otherwise.
func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// isLetterOrDigit returns true if r is a Unicode letter, an underscore or a
// digit, and false otherwise.
func isLetterOrDigit(r rune) bool {
	return isLetter(r) || unicode.IsDigit(r)
}

// isDigit returns true if r is an ASCII digit from '0' through '9', and false
// otherwise.
func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

// isSpace returns true if r is a whitespace character, and false otherwise.
func isSpace(r rune) bool {
	// whitespace characters, except newline.
	const whitespace = " \t"

	return strings.IndexRune(whitespace, r) != -1
}

// accept consumes the next rune if it's from the valid set. It returns true if
// a rune was consumed and false otherwise.
func (l *lexer) accept(valid string) bool {
	r := l.next()
	if r == eof {
		return false
	}
	if strings.IndexRune(valid, r) == -1 {
		l.backup()
		return false
	}
	return true
}

// acceptFunc consumes the next rune if it's valid. It returns true if a rune
// was consumed and false otherwise.
func (l *lexer) acceptFunc(isValid func(rune) bool) bool {
	r := l.next()
	if r == eof {
		return false
	}
	if !isValid(r) {
		l.backup()
		return false
	}
	return true
}

// acceptRun consumes a run of runes from the valid set. It returns true if a
// rune was consumed and false otherwise.
func (l *lexer) acceptRun(valid string) bool {
	consumed := false
	for l.accept(valid) {
		consumed = true
	}
	return consumed
}

// acceptRunFunc consumes a run of valid runes. It returns true if a rune was
// consumed and false otherwise.
func (l *lexer) acceptRunFunc(isValid func(rune) bool) bool {
	consumed := false
	for l.acceptFunc(isValid) {
		consumed = true
	}
	return consumed
}

// acceptIdent consumes an identifier. It returns true if a rune was consumed
// and false otherwise.
func (l *lexer) acceptIdent() bool {
	if !l.acceptFunc(isLetter) {
		return false
	}
	l.acceptRunFunc(isLetterOrDigit)
	return true
}

// ignore ignores any pending input read since the last token.
func (l *lexer) ignore() {
	l.start = l.pos
}

// ignoreRun ignores a run of valid runes.
func (l *lexer) ignoreRun(valid string) {
	if l.acceptRun(valid) {
		l.ignore()
	}
}

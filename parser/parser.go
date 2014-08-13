// Package parser implements syntactical parsing of assembly source code; the
// output of which is an abstract syntax tree (AST).
package parser

import (
	"log"

	"github.com/mewmew/asm/ast"
	"github.com/mewmew/asm/token"
)

// A parser parses the input tokens into an AST.
type parser struct {
	// The input tokens.
	input []token.Token
	// The current position in slice of the input tokens.
	pos int
	// A slice of parsed AST nodes.
	nodes []ast.Node
}

// Parse parses the provided slice of assembly tokens into an AST.
func Parse(tokens []token.Token) (nodes []ast.Node) {
	p := &parser{
		input: tokens,
		nodes: make([]ast.Node, 0),
	}

	// Parse the input tokens.
	p.run()

	return p.nodes
}

// run parses the input by repeatedly executing the active state function until
// it returns a nil state.
func (p *parser) run() {
	// parseLine is the initial state function of the parser.
	for state := parseLine; state != nil; {
		state = state(p)
	}
}

// next consumes the next token in the slice of input tokens.
func (p *parser) next() (tok token.Token) {
	if p.pos >= len(p.input) {
		log.Fatalln("parser.next: TODO")
	}
	tok = p.input[p.pos]
	p.pos++
	return tok
}

// peek returns but does not consume the next token in the slice of the input
// tokens.
func (p *parser) peek() (tok token.Token) {
	tok = p.next()
	if tok.Typ == token.EOF {
		return tok
	}
	p.backup()
	return tok
}

// backup backs up one token in the slice of input tokens.
func (p *parser) backup() {
	if p.pos == 0 {
		log.Fatalln("parser.backup: invalid call; no tokens have yet been consumed.")
	}
	p.pos--
}

// Package token defines constants representing the lexical tokens of the
// assembly language.
package token

import "fmt"

// A Token represents a lexical token of the assembly language.
type Token struct {
	// The token type.
	Typ Type
	// The string value of the token.
	Val string
}

func (tok Token) String() string {
	if len(tok.Val) == 0 {
		return tok.Typ.String()
	}
	return fmt.Sprintf("[%v]: %q", tok.Typ, tok.Val)
}

// Type is the set of lexical token types of the assembly language.
type Type int

// Token types.
const (
	Error       Type = iota // an error occurred; value contains the error message.
	EOF                     // end of file.
	Ident                   // identifier.
	Colon                   // label declaration colon.
	Int                     // integer literal
	Char                    // character literal
	String                  // string literal
	Comma                   // operand separation comma
	LineComment             // ; line comment
	Newline                 // newline
)

// names is a map from token type to token name.
var names = map[Type]string{
	Error:       "error",
	EOF:         "EOF",
	Ident:       "identifier",
	Colon:       ":",
	Int:         "integer literal",
	Char:        "character literal",
	String:      "string literal",
	Comma:       ",",
	LineComment: "line comment",
	Newline:     "newline",
}

func (typ Type) String() string {
	name, ok := names[typ]
	if !ok {
		return fmt.Sprintf("<unknown token type: %d>", typ)
	}
	return name
}

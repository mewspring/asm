package lexer

import (
	"errors"
	"fmt"

	"github.com/mewmew/asm/token"
)

// A stateFn represents the state of the scanner as a function that returns a
// state function.
type stateFn func(l *lexer) stateFn

// lexLine lexes a line. It is the initial state function of the lexer.
func lexLine(l *lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			return l.emitEOF()
		case isSpace(r):
			// Ignore whitespace.
			l.ignore()
		case isLetter(r):
			l.backup()
			// Lex label, operator, register or directive identifier.
			return lexIdent
		case r == '.':
			// Lex local label identifier.
			// TODO(u): Add support for floating-point literals.
			return lexIdent
		case r == ':':
			// Lex label declaration colon.
			return lexColon
		case r == '+':
			if isDigit(l.peek()) {
				return lexIntLit
			}
			return lexAddOp
		case r == '-':
			if isDigit(l.peek()) {
				return lexIntLit
			}
			return lexSubOp
		case r == '*':
			return lexMulOp
		case r == '/':
			return lexDivOp
		case r == '%':
			return lexModOp
		case r == '&':
			return lexAndOp
		case r == '|':
			return lexOrOp
		case r == '^':
			return lexXorOp
		case r == '<':
			return lexLshOp
		case r == '>':
			return lexRshOp
		case isDigit(r):
			// Lex integer literal.
			l.backup()
			return lexIntLit
		case r == '\'':
			// Lex character literal.
			return lexCharLit
		case r == '"':
			// Lex string literal.
			return lexStringLit
		case r == '`':
			// Lex raw string literal.
			return lexRawStringLit
		case r == '$':
			// Lex register identifier.
			return lexIdent
		case r == ',':
			// Lex operand separation comma
			return lexComma
		case r == ';':
			// Lex line comment.
			return lexLineComment
		case r == '\n':
			// Lex newline.
			return lexNewline
		default:
			return l.errorf("lexer.lexLine: unexpected rune '%c' at beginning of token", r)
		}
	}
}

// lexIdent lexes an identifier. A '.' or a '$' may have been consumed already.
func lexIdent(l *lexer) stateFn {
	if !l.acceptFunc(isLetter) {
		return l.errorf("lexer.lexIdent: expected Unicode letter or underscore, got '%c'", l.next())
	}
	l.acceptRunFunc(isLetterOrDigit)
	l.emit(token.Ident)
	return lexLine
}

// lexColon lexes a label declaration colon. A ':' has already been consumed.
func lexColon(l *lexer) stateFn {
	l.emit(token.Colon)
	return lexLine
}

// lexIntLit lexes an integer literal. No characters have been consumed already.
func lexIntLit(l *lexer) stateFn {
	// Optional sign.
	l.accept("+-")

	// Decimal digits.
	digits := "0123456789"
	hasIntLit := false
	r := l.next()
	switch {
	case r == eof:
		return l.emitEOF()
	case r == '0':
		switch {
		case l.accept("xX"):
			// Hexadecimal digits.
			digits = "0123456789abcdefABCDEF"
		case l.accept("b"):
			// Binary digits.
			digits = "01"
		default:
			// Octal digits.
			digits = "01234567"

			// The octal literal '0' has been consumed.
			hasIntLit = true
		}
	case r >= '1' && r <= '9':
		// A decimal literal '1' … '9' has been consumed.
		hasIntLit = true
	default:
		return l.errorf("lexer.lexIntLit: expected decimal digit, got %c", r)
	}
	if !l.acceptRun(digits) && !hasIntLit {
		return l.errorf("lexer.lexIntLit: missing digits in integer literal")
	}
	l.emit(token.Int)

	return lexLine
}

// lexCharLit lexes a character literal. A single quote has already been
// consumed.
func lexCharLit(l *lexer) stateFn {
	if l.accept(`\`) {
		// Consume backslash escape sequence.
		err := consumeEscape(l, '\'')
		if err != nil {
			return l.errorf("lexer.lexCharLit: %v", err)
		}
	} else {
		// Consume single quoted character. Any character may appear except single
		// quote and newline.
		r := l.next()
		switch r {
		case eof:
			return l.errorf("lexer.lexCharLit: unexpected eof in single quoted character literal")
		case '\'':
			return l.errorf("lexer.lexCharLit: unexpected ' in character literal")
		case '\n':
			return l.errorf("lexer.lexCharLit: unexpected newline in character literal")
		}
	}
	if !l.accept("'") {
		return l.errorf("lexer.lexCharLit: missing ' in character literal")
	}
	l.emit(token.Char)
	return lexLine
}

// lexStringLit lexes a string literal. A '"' has already been consumed.
func lexStringLit(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			return l.errorf("lexer.lexStringLiteral: unexpected eof in string literal")
		case '\n':
			return l.errorf("lexer.lexStringLiteral: unexpected newline in string literal")
		case '\\':
			// Consume backslash escape sequence.
			err := consumeEscape(l, '"')
			if err != nil {
				return l.errorf("lexer.lexStringLit: %v", err)
			}
		case '"':
			l.emit(token.String)
			return lexLine
		}
	}
}

// lexRawStringLit lexes a raw string literal. A '`' has already been consumed.
func lexRawStringLit(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			return l.errorf("lexer.lexRawStringLiteral: unexpected eof in raw string literal")
		case '`':
			l.emit(token.String)
			return lexLine
		}
	}
}

// consumeEscape consumes an escape sequence. A backslash has already been
// consumed. A valid single-character escape sequence is specified by valid.
// Single quotes are only valid within character literals and double quotes are
// only valid within string literals.
func consumeEscape(l *lexer, valid rune) (err error) {
	// Several backslash escapes allow arbitrary values to be encoded as ASCII
	// text. There are two ways to represent the integer value as a numeric
	// constant: \x followed by exactly two hexadecimal digits, and a plain
	// backslash \ followed by exactly three octal digits. Although these
	// representations all result in an integer, they have different valid
	// ranges. Octal escapes must represent a value between 0 and 255 inclusive.
	//
	// After a backslash, certain single-character escapes represent special
	// values:
	//
	//    \a   U+0007 alert or bell
	//    \b   U+0008 backspace
	//    \f   U+000C form feed
	//    \n   U+000A line feed or newline
	//    \r   U+000D carriage return
	//    \t   U+0009 horizontal tab
	//    \v   U+000b vertical tab
	//    \\   U+005c backslash
	//    \'   U+0027 single quote  (valid escape only within character literals)
	//    \"   U+0022 double quote  (valid escape only within string literals)
	//
	// All other sequences starting with a backslash are illegal inside character
	// literals.
	r := l.next()
	switch r {
	case eof:
		return errors.New("lexer.consumeEscape: unexpected eof after backslash escape character")
	case '0', '1', '2', '3':
		// Octal escape sequence.
		digits := "01234567"
		if !l.accept(digits) || !l.accept(digits) {
			return fmt.Errorf("lexer.consumeEscape: non-octal character in escape sequence: %c", l.next())
		}
	case 'x':
		// Hexadecimal escape sequence.
		digits := "0123456789abcdefABCDEF"
		if !l.accept(digits) || !l.accept(digits) {
			return fmt.Errorf("lexer.consumeEscape: non-hex character in escape sequence: %c", l.next())
		}
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\':
		// Consume single-character escape.
	default:
		// Single quote escapes are only valid within character literals and
		// double quote escapes are only valid within string literals.
		if r != valid {
			return fmt.Errorf("lexer.consumeEscape: unknown escape sequence: %c", r)
		}
	}

	return nil
}

// lexComma lexes an operand separation comma. A ',' has already been consumed.
func lexComma(l *lexer) stateFn {
	l.emit(token.Comma)
	return lexLine
}

// lexLineComment lexes a line comment. A ';' has already been consumed.
func lexLineComment(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			l.emit(token.LineComment)
			l.emit(token.EOF)
			return nil
		case '\n':
			l.backup()
			l.emit(token.LineComment)
			return lexLine
		}
	}
}

// lexNewline lexes a newline. A '\n' has already been consumed.
func lexNewline(l *lexer) stateFn {
	l.emit(token.Newline)
	return lexLine
}

func lexAddOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexSubOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexMulOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexDivOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexModOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexAndOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexOrOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexXorOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexLshOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

func lexRshOp(l *lexer) stateFn {
	panic("not yet implemented.")
}

package lexer

import (
	"io/ioutil"
	"testing"

	"github.com/mewlang/asm/token"
)

type test struct {
	path   string
	tokens []token.Token
}

func TestParse(t *testing.T) {
	golden := []test{
		// i=0
		{
			path: "../testdata/hello_x86.asm",
			tokens: []token.Token{
				token.Token{Typ: token.LineComment, Val: "; === [ text section ] ========================================================="},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "section"},
				token.Token{Typ: token.String, Val: "\".text\""},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "global"},
				token.Token{Typ: token.Ident, Val: "_start"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "_start"},
				token.Token{Typ: token.Colon, Val: ":"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.LineComment, Val: "; write(1, \"Hello world!\\n\", 13)"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "mov"},
				token.Token{Typ: token.Ident, Val: "eax"},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Int, Val: "4"},
				token.Token{Typ: token.LineComment, Val: "; sys_write"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "mov"},
				token.Token{Typ: token.Ident, Val: "ebx"},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Int, Val: "1"},
				token.Token{Typ: token.LineComment, Val: ";    fd: stdout"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "mov"},
				token.Token{Typ: token.Ident, Val: "ecx"},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Ident, Val: "hello"},
				token.Token{Typ: token.LineComment, Val: ";    buf: str"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "mov"},
				token.Token{Typ: token.Ident, Val: "edx"},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Int, Val: "13"},
				token.Token{Typ: token.LineComment, Val: ";    len: len(str)"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "int"},
				token.Token{Typ: token.Int, Val: "0x80"},
				token.Token{Typ: token.LineComment, Val: "; syscall"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.LineComment, Val: "; exit(0)"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "mov"},
				token.Token{Typ: token.Ident, Val: "eax"},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Int, Val: "1"},
				token.Token{Typ: token.LineComment, Val: "; sys_exit"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "mov"},
				token.Token{Typ: token.Ident, Val: "ebx"},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Int, Val: "0"},
				token.Token{Typ: token.LineComment, Val: ";    status = 0"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "int"},
				token.Token{Typ: token.Int, Val: "0x80"},
				token.Token{Typ: token.LineComment, Val: "; syscall"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.LineComment, Val: "; === [ rdata section ] ========================================================"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "section"},
				token.Token{Typ: token.String, Val: "\".rdata\""},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.Ident, Val: "hello"},
				token.Token{Typ: token.Colon, Val: ":"},
				token.Token{Typ: token.Ident, Val: "db"},
				token.Token{Typ: token.String, Val: "\"Hello world\""},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Char, Val: "'!'"},
				token.Token{Typ: token.Comma, Val: ","},
				token.Token{Typ: token.Int, Val: "10"},
				token.Token{Typ: token.Newline, Val: "\n"},
				token.Token{Typ: token.EOF, Val: ""},
			},
		},
	}

	for i, g := range golden {
		buf, err := ioutil.ReadFile(g.path)
		if err != nil {
			t.Error(err)
		}
		input := string(buf)

		tokens := Parse(input)
		for j, want := range g.tokens {
			if len(tokens) < j {
				t.Errorf("missing token; unable to access token at index %d", j)
			}
			got := tokens[j]
			if got != want {
				t.Errorf("i=%d: expected %v, got %v", i, want, got)
			}
		}
	}
}

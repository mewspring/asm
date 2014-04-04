package lexer_test

import (
	"fmt"

	"github.com/mewmew/asm/lexer"
)

func ExampleParse() {
	input := `; simple loop in MIPS assembly.
	ori     $t0, $zero, 16
loop:
	beq     $t0, $zero done
	addi    $t0, $t0, -1
	; do stuff.
	j       loop
done:
`
	tokens := lexer.Parse(input)
	for i, tok := range tokens {
		fmt.Printf("token %d: %v\n", i, tok)
	}

	// Output:
	// token 0: [line comment]: "; simple loop in MIPS assembly."
	// token 1: [newline]: "\n"
	// token 2: [identifier]: "ori"
	// token 3: [identifier]: "$t0"
	// token 4: [,]: ","
	// token 5: [identifier]: "$zero"
	// token 6: [,]: ","
	// token 7: [integer literal]: "16"
	// token 8: [newline]: "\n"
	// token 9: [identifier]: "loop"
	// token 10: [:]: ":"
	// token 11: [newline]: "\n"
	// token 12: [identifier]: "beq"
	// token 13: [identifier]: "$t0"
	// token 14: [,]: ","
	// token 15: [identifier]: "$zero"
	// token 16: [identifier]: "done"
	// token 17: [newline]: "\n"
	// token 18: [identifier]: "addi"
	// token 19: [identifier]: "$t0"
	// token 20: [,]: ","
	// token 21: [identifier]: "$t0"
	// token 22: [,]: ","
	// token 23: [integer literal]: "-1"
	// token 24: [newline]: "\n"
	// token 25: [line comment]: "; do stuff."
	// token 26: [newline]: "\n"
	// token 27: [identifier]: "j"
	// token 28: [identifier]: "loop"
	// token 29: [newline]: "\n"
	// token 30: [identifier]: "done"
	// token 31: [:]: ":"
	// token 32: [newline]: "\n"
	// token 33: EOF
}

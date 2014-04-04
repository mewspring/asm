// tokens demonstrates how to tokenize input files using the Parse function.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mewmew/asm/lexer"
)

func main() {
	flag.Parse()
	for _, filePath := range flag.Args() {
		err := tokens(filePath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// tokens lexes the provided file and prints each token to standard output.
func tokens(filePath string) (err error) {
	// Read input from file.
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	input := string(buf)

	// Tokenize the input.
	tokens := lexer.Parse(input)
	for _, tok := range tokens {
		fmt.Println("tok:", tok)
	}

	return nil
}

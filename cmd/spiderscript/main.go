package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/chronos-tachyon/go-spiderscript/ast"
	"github.com/chronos-tachyon/go-spiderscript/token"
)

func main() {
	flag.Parse()

	for _, inputFile := range flag.Args() {
		bytes, err := ioutil.ReadFile(inputFile)
		if err != nil {
			panic(err)
		}

		runes := make([]rune, 0, len(bytes))
		offset := uint(0)
		for offset < uint(len(bytes)) {
			ch, n := utf8.DecodeRune(bytes[offset:])
			if n < 1 {
				n = 1
			}
			runes = append(runes, ch)
			offset += uint(n)
		}

		var lexer token.Lexer
		lexer.Init(inputFile, runes)

		var parser ast.Parser
		parser.Init(&lexer)

		var file ast.File
		parser.Parse(&file)
		fmt.Printf("%v\n", &file)

		if errors := parser.Errors(); errors == nil {
			fmt.Fprintf(os.Stderr, "No errors.\n")
		} else {
			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}
	}
}

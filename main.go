package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KyAnhVo/goviz/lexer"
)

func main() {
	lexerCheck()
}

func lexerCheck() {
	inputFile := flag.String("from", "", "File to lex from")
	outputFile := flag.String("to", "stdout", "to: 'stdout' | file path")
	flag.Parse()

	var output io.Writer
	var file *os.File
	var err error

	if *outputFile == "stdout" {
		output = os.Stdout
	} else {
		file, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		output = file
	}

	if *inputFile == "" {
		fmt.Println("No input file stated")
		return
	}
	content, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	l := lexer.NewLexer([]rune(string(content)))

	var token lexer.Token
	var pos lexer.Pos
	token, pos, err = l.GetNextToken()
	for pos != lexer.PosEOF {
		output.Write(fmt.Appendf(
			[]byte(""), "Token: %s\nPos: %+v\n\n", lexer.FormatToken(token), pos,
		))
		token, pos, err = l.GetNextToken()
	}
}

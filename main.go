package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KyAnhVo/goviz/lexer"
	types "github.com/KyAnhVo/goviz/token"
)

func main() {
	lexerCheck()
}

func lexerCheck() {
	inputFile := flag.String("from", "", "File to lex from")
	outputFile := flag.String("to", "stdout", "to: 'stdout' | file path")
	flag.Parse()

	var output io.Writer
	var out, in *os.File
	var err error

	if *outputFile == "stdout" {
		output = os.Stdout
	} else {
		out, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println("Error creating file:\n", err)
			return
		}
		defer out.Close()

		output = out
	}

	if *inputFile == "" {
		fmt.Println("No input file stated")
		return
	} else {
		in, err = os.Open(*inputFile)
		if err != nil {
			fmt.Println("Error reading file:\n", err)
			return
		}
		defer in.Close()

	}

	l, err := lexer.NewLexer(*bufio.NewScanner(in))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return
	}

	var token types.Token
	var pos types.Pos
	token, pos, err = l.GetNextToken()
	for pos != types.PosEOF {
		output.Write(fmt.Appendf(
			[]byte(""), "Token: %s\nPos: %+v\n\n", types.FormatToken(token), pos,
		))
		token, pos, err = l.GetNextToken()
	}
}

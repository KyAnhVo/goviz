package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/KyAnhVo/goviz/lexer"
	"github.com/KyAnhVo/goviz/parser"
	"github.com/KyAnhVo/goviz/token"
)

const (
	noneMode   string = "none"
	lexerMode  string = "l"
	parserMode string = "p"
)

func main() {
	inputFile := flag.String("from", "", "File to lex from")
	outputFile := flag.String("to", "stdout", "to: 'stdout' | file path")
	modeFlag := flag.String("mode", "none", "mode: 'l' | 'p' | 'none'")

	// flag processing
	flag.Parse()
	if inputFile == nil {
		io.WriteString(os.Stderr, "Must have input file")
		return
	}
	if *modeFlag != noneMode && *modeFlag != lexerMode && *modeFlag != parserMode {
		log.Fatalf("Mode flag must be %s, %s, or %s", noneMode, lexerMode, parserMode)
		return
	}

	var output *os.File
	var err error
	if *outputFile == "stdout" {
		output = os.Stdout
	} else {
		output, err = os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Lexer create error: %s", err.Error())
			os.Exit(1)
		}
	}

	l, err := createLexer(*inputFile)
	if err != nil {
		log.Fatalf("Lexer create error: %s", err.Error())
		os.Exit(1)
	}
	p, err := parser.New(l)
	if err != nil {
		log.Fatalf("Parser create error: %s", err.Error())
		os.Exit(1)
	}

	if *modeFlag == lexerMode {
		runLexer(l, output)
	}

	if *modeFlag == parserMode {
		runParser(p, output)
	}
}

func createLexer(inputFile string) (*lexer.Lexer, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	return lexer.NewLexer(scanner)
}

func runLexer(l *lexer.Lexer, output *os.File) {
	var tok token.Token
	var err error
	var pos token.Pos

	tok, pos, err = l.GetNextToken()
	for pos != token.PosEOF {
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
			os.Exit(1)
		}
		io.WriteString(output,
			fmt.Sprintf("%s\n%s\n\n",
				token.FormatToken(tok), token.FormatPos(pos)))
		tok, pos, err = l.GetNextToken()
	}
}

func runParser(p *parser.Parser, output *os.File) {
	ast, err := p.Parse()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
		os.Exit(1)
	}
	io.WriteString(output, fmt.Sprintf(
		"%+v", ast))
}

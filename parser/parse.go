package parser

import (
	"github.com/KyAnhVo/goviz/ast"
	"github.com/KyAnhVo/goviz/token"
)

func (p *Parser) Parse() (ast.AST, error) {
	ast := ast.AST{}
	var err error

	ast.PackageName, err = p.packageDecl()
	if err != nil {
		return ast, err
	}

	ast.ImportedFiles, err = p.importDecls()
	if err != nil {
		return ast, err
	}

	/*
		for topLevelDeclStarters.Contains(p.currToken.token) {
			p.topLevelDecl()
		}
	*/

	return ast, nil
}

func (p *Parser) packageDecl() (string, error) {
	var err error

	err = p.ExpectToken(token.TokenKeywordPackage)
	if err != nil {
		return "", err
	}
	p.advance()

	err = p.ExpectTokenType(token.TokenTypeIdentifier)
	if err != nil {
		return "", err
	}
	name := p.currToken.token.Value
	p.advance()

	err = p.ExpectToken(token.TokenSemicolon)
	if err != nil {
		return "", err
	}
	p.advance()

	return name, nil
}

// ImportDecls = { "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) } ";" .
func (p *Parser) importDecls() (ast.Imports, error) {
	imports := make(ast.Imports, 0, 8)

	for p.currToken.token == token.TokenKeywordImport {
		p.advance()
		if p.currToken.token == token.TokenLParen {
			p.advance()
			for p.currToken.token != token.TokenRParen {
				importSpec, err := p.importSpec()
				if err != nil {
					return []ast.Import{}, err
				}
				imports = append(imports, importSpec)

				err = p.ExpectToken(token.TokenSemicolon)
				if err != nil {
					err2 := p.ExpectToken(token.TokenRParen)
					if err2 != nil {
						// problem is forgetting to insert semicolon,
						// not the not-closing
						return []ast.Import{}, err
					} else {
						break
					}
				}

				p.advance()
			}
			p.advance()
		} else {
			// importSpec path
			importSpec, err := p.importSpec()
			if err != nil {
				return []ast.Import{}, err
			}
			imports = append(imports, importSpec)

			err = p.ExpectToken(token.TokenSemicolon)
			if err != nil {
				return []ast.Import{}, err
			}
			p.advance()
		}

		err := p.ExpectToken(token.TokenSemicolon)
		if err != nil {
			return []ast.Import{}, err
		}
		p.advance()
	}
	return imports, nil
}

// ImportSpec = [ "." | PackageName ] ImportPath .
// ImportPath = string_lit .
func (p *Parser) importSpec() (ast.Import, error) {
	var importVal ast.Import
	tok := p.currToken.token
	if tok == token.TokenDot {
		importVal.Name = "."
		p.advance()
	} else if tok.Type == token.TokenTypeIdentifier {
		importVal.Name = tok.Value
		p.advance()
	} else {
		importVal.Name = ""
	}

	p.ExpectTokenType(token.TokenTypeStringLiteral)
	importVal.Src = tok.Value
	p.advance()
	return importVal, nil
}

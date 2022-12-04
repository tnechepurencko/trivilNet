package ast

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

type Parser struct {
	lex    *lexer.Lexer
	module *ast.Module
}

func (p *Parser) Init(source *env.Source) {

}

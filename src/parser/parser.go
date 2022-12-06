package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

type Parser struct {
	source *env.Source
	lex    *lexer.Lexer
	module *ast.Module

	pos     int
	tok     lexer.Token
	lit     string
	afterNL bool

	trace  bool
	indent int
}

func Parse(source *env.Source) *ast.Module {
	p := new(Parser)
	p.lex = new(lexer.Lexer)
	p.lex.Init(source)
	p.source = source

	p.trace = true //!

	p.next()
	p.parseModule()

	return p.module
}

func (p *Parser) next() {
	p.afterNL = false
	for true {
		p.pos, p.tok, p.lit = p.lex.Scan()

		switch p.tok {
		case lexer.EOF:
			p.afterNL = true
			return
		case lexer.BLOCK_COMMENT:
			continue
		case lexer.LINE_COMMENT:
			p.afterNL = true
			continue
		case lexer.NL:
			p.afterNL = true
			continue
		default:
			return
		}
	}
}

func (p *Parser) expect(tok lexer.Token) {
	if p.tok != tok {
		env.AddError(p.pos, "ПАР-ОЖИДАЛСЯ", tok.String())
	}
	p.next()
}

//=====

func (p *Parser) parseModule() {

	p.module = ast.NewModule()

	if p.tok != lexer.MODULE {
		env.AddError(p.pos, "ПАР-ОЖИДАЛСЯ", lexer.MODULE.String())
		return
	}
	p.next()
	p.module.Name = p.parseIdent()

	p.parseImportList()
	p.parseDeclarations()
}

func (p *Parser) parseImportList() {
	if p.trace {
		defer un(trace(p, "Import List"))
	}

}

func (p *Parser) parseDeclarations() {
	if p.trace {
		defer un(trace(p, "Declarations"))
	}

}

//====

func (p *Parser) parseIdent() string {
	name := "_"
	if p.tok == lexer.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(lexer.IDENT)
	}
	return name
}

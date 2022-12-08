package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

//=== statements

func (p *Parser) parseStatementSeq() *ast.StatementSeq {

	var n = &ast.StatementSeq{
		StatementBase: ast.StatementBase{Pos: p.pos},
		Statements:    make([]ast.Statement, 0),
	}

	if p.tok != lexer.LBRACE {
		p.expect(lexer.LBRACE)
		return n
	}
	p.next()
	for p.tok != lexer.EOF && p.tok != lexer.RBRACE {
		var s = p.parseStatement()
		if s != nil {
			n.Statements = append(n.Statements, s)
		}
		!sep
	}

	p.expect(lexer.RBRACE)

	return n
}

func (p *Parser) parseStatement() ast.Statement {

	switch p.tok {
	default:
		env.AddError(p.pos, "ПАР-ОШ-ОПЕРАТОР", p.tok.String())
		return nil
	}

	return nil
}

//====

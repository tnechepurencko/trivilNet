package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

var validSimpleStmToken = map[lexer.Token]bool{
	lexer.IDENT: true,
	lexer.LPAR:  true,

	// literals
	lexer.INT:    true,
	lexer.FLOAT:  true,
	lexer.STRING: true,

	// unary ops
	lexer.ADD: true,
	lexer.SUB: true,
	lexer.NOT: true,
}

//=== statements

func (p *Parser) parseStatementSeq() *ast.StatementSeq {
	if p.trace {
		defer un(trace(p, "Список операторов"))
	}

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

		if p.tok == lexer.RBRACE {
			break
		}
		p.expectSep("ПАР_РАЗД_ОПЕРАТОРОВ")
	}

	p.expect(lexer.RBRACE)

	return n
}

func (p *Parser) parseStatement() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор"))
	}

	switch p.tok {
	default:
		if validSimpleStmToken[p.tok] {
			return p.parseSimpleStatement()
		}
		env.AddError(p.pos, "ПАР-ОШ-ОПЕРАТОР", p.tok.String())
		return nil
	}

	return nil
}

func (p *Parser) parseSimpleStatement() ast.Statement {
	if p.trace {
		defer un(trace(p, "Простой оператор"))
	}

	var expr = p.parseExpression()

	var s = &ast.ExprStatement{
		StatementBase: ast.StatementBase{Pos: p.pos},

		X: expr,
	}

	return s
}

//====

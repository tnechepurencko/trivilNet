package parser

import (
	"fmt"
	"trivil/ast"

	//	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

//=== приоритеты

const lowestPrecedence = 0

func precedence(tok lexer.Token) int {
	switch tok {
	case lexer.OR:
		return 1
	case lexer.AND:
		return 2
	case lexer.EQ, lexer.NEQ, lexer.LSS, lexer.LEQ, lexer.GTR, lexer.GEQ:
		return 3
	case lexer.ADD, lexer.SUB, lexer.BITOR:
		return 4
	case lexer.MUL, lexer.QUO, lexer.REM, lexer.BITAND:
		return 5
	default:
		return lowestPrecedence
	}
}

//=== выражения

func (p *Parser) parseExpression() ast.Expr {
	if p.trace {
		defer un(trace(p, "Выражение"))
	}

	return p.parseBinaryExpression(lowestPrecedence + 1)
}

func (p *Parser) parseBinaryExpression(prec int) ast.Expr {
	if p.trace {
		defer un(trace(p, "Выражение бинарное"))
	}

	var x = p.parseUnaryExpression()
	for {
		op := p.tok
		opPrec := precedence(op)
		if opPrec < prec {
			return x
		}
		var pos = p.pos
		p.next()

		var y = p.parseBinaryExpression(opPrec + 1)
		x = &ast.BinaryExpr{
			ExprBase: ast.ExprBase{Pos: pos},
			X:        x,
			Op:       op,
			Y:        y,
		}
	}
}

func (p *Parser) parseUnaryExpression() ast.Expr {
	if p.trace {
		defer un(trace(p, "Выражение унарное"))
	}

	switch p.tok {
	case lexer.SUB, lexer.NOT:
		var pos = p.pos
		var op = p.tok
		p.next()
		var x = p.parseUnaryExpression()
		return &ast.UnaryExpr{
			ExprBase: ast.ExprBase{Pos: pos},
			Op:       op,
			X:        x,
		}
	case lexer.ADD:
		return p.parseUnaryExpression()
	}

	var x = p.parsePrimaryExpression()

	// check ?
	return x
}

func (p *Parser) parsePrimaryExpression() ast.Expr {
	if p.trace {
		defer un(trace(p, "Выражение первичное"))
	}

	return nil
}

package parser

import (
	"fmt"
	"trivil/ast"
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

	var x ast.Expr

	switch p.tok {
	case lexer.INT, lexer.FLOAT, lexer.STRING:
		x = &ast.LiteralExpr{
			ExprBase: ast.ExprBase{Pos: p.pos},
			Kind:     p.tok,
			Lit:      p.lit,
		}
		p.next()
	case lexer.IDENT:
		x = &ast.IdentExpr{
			ExprBase: ast.ExprBase{Pos: p.pos},
			Name:     p.lit,
		}
		p.next()
	case lexer.LPAR:
		p.next()
		x = p.parseExpression()
		p.expect(lexer.RPAR)
	default:
		p.error(p.pos, "ПАР-ОШ-ОПЕРАНД", p.tok.String())
		return &ast.InvalidExpr{}
	}

	for {
		switch p.tok {
		case lexer.DOT:
			x = p.parseSelector(x)
		case lexer.LPAR:
			x = p.parseArguments(x)
		case lexer.LCONV:
			x = p.parseConversion(x)
		case lexer.LBRACK:
			x = p.parseIndex(x)
		default:
			return x
		}
	}

}

func (p *Parser) parseSelector(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Селектор"))
	}

	var n = &ast.SelectorExpr{
		ExprBase: ast.ExprBase{Pos: p.pos},
		X:        x,
	}

	p.next()
	n.Name = p.parseIdent()

	return n
}

func (p *Parser) parseArguments(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Аргументы"))
	}

	var n = &ast.CallExpr{
		ExprBase: ast.ExprBase{Pos: p.pos},
		X:        x,
		Args:     make([]ast.Expr, 0),
	}

	p.expect(lexer.LPAR)

	for p.tok != lexer.RPAR && p.tok != lexer.EOF {

		var expr = p.parseExpression()

		n.Args = append(n.Args, expr)

		if p.tok == lexer.RPAR {
			break
		}
		p.expect(lexer.COMMA)
	}

	p.expect(lexer.RPAR)

	return n
}

func (p *Parser) parseConversion(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Конверсия"))
	}

	var n = &ast.ConversionExpr{
		ExprBase: ast.ExprBase{Pos: p.pos},
		X:        x,
	}

	p.next()
	n.Typ = p.parseTypeRef()

	p.expect(lexer.RCONV)

	return n
}

func (p *Parser) parseIndex(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Индексация"))
	}

	var n = &ast.IndexExpr{
		ExprBase: ast.ExprBase{Pos: p.pos},
		X:        x,
		Values:   make([]ast.ValuePair, 0),
	}

	p.expect(lexer.LBRACK)

	var l ast.Expr
	var r ast.Expr

	for p.tok != lexer.RBRACK && p.tok != lexer.EOF {

		l = p.parseExpression()

		if p.tok == lexer.COLON {
			p.next()
			r = p.parseExpression()
		} else {
			r = nil
		}

		n.Values = append(n.Values, ast.ValuePair{L: l, R: r})

		if p.tok == lexer.RBRACK {
			break
		}
		p.expect(lexer.COMMA)
	}

	p.expect(lexer.RBRACK)

	p.checkIndexValues(n)

	return n
}

func (p *Parser) checkIndexValues(n *ast.IndexExpr) {

	var pairs = 0
	for _, v := range n.Values {
		if v.R != nil {
			pairs++
		}
	}

	if pairs == len(n.Values) {
		n.Pairs = true
	} else if pairs != 0 {
		p.error(n.Pos, "ПАР-СМЕСЬ-МАССИВ")
	}
}

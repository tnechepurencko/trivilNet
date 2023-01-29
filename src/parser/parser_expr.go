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
	case lexer.INT:
		x = &ast.LiteralExpr{
			ExprBase: ast.ExprBase{Pos: p.pos},
			Kind:     ast.Lit_Int,
			Lit:      p.lit,
		}
		p.next()
	case lexer.FLOAT:
		x = &ast.LiteralExpr{
			ExprBase: ast.ExprBase{Pos: p.pos},
			Kind:     ast.Lit_Float,
			Lit:      p.lit,
		}
		p.next()
	case lexer.STRING:
		x = &ast.LiteralExpr{
			ExprBase: ast.ExprBase{Pos: p.pos},
			Kind:     ast.Lit_String,
			Lit:      p.lit,
		}
		p.next()
	case lexer.SYMBOL:
		x = &ast.LiteralExpr{
			ExprBase: ast.ExprBase{Pos: p.pos},
			Kind:     ast.Lit_Symbol,
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
		case lexer.LBRACE:
			if p.lex.WhitespaceBefore('{') {
				return x
			}
			x = p.parseClassComposite(x)

		default:
			return x
		}
	}

}

func (p *Parser) parseSelector(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Селектор"))
	}
	p.next()

	var n = &ast.SelectorExpr{
		ExprBase: ast.ExprBase{Pos: p.pos},
		X:        x,
	}

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

		if p.tok == lexer.ELLIPSIS {
			var u = &ast.UnfoldExpr{
				ExprBase: ast.ExprBase{Pos: p.pos},
				X:        expr,
			}
			expr = u
			p.next()
		}

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

	if p.tok == lexer.CAUTION {
		p.next()
		n.Caution = true
		if !p.module.Caution {
			p.error(p.pos, "ПАР-ОШ-ИСП-ОСТОРОЖНО")
		}
	}

	n.TargetTyp = p.parseTypeRef()

	p.expect(lexer.RPAR)

	return n
}

func (p *Parser) parseIndex(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Индексация"))
	}

	var n = &ast.GeneralBracketExpr{
		ExprBase: ast.ExprBase{Pos: p.pos},
		X:        x,
		Composite: &ast.ArrayCompositeExpr{
			ExprBase: ast.ExprBase{Pos: p.pos},
			Elements: make([]ast.ElementPair, 0),
		},
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
			r = l
			l = nil
		}

		n.Composite.Elements = append(n.Composite.Elements, ast.ElementPair{Key: l, Value: r})

		if p.tok == lexer.RBRACK {
			break
		}
		p.expect(lexer.COMMA)
	}

	p.expect(lexer.RBRACK)

	p.checkElements(n.Composite)

	return n
}

func (p *Parser) checkElements(n *ast.ArrayCompositeExpr) {

	var pairs = 0
	for _, v := range n.Elements {
		if v.Key != nil {
			pairs++
		}
	}

	if pairs == len(n.Elements) {
		n.Keys = true
	} else if pairs != 0 {
		p.error(n.Pos, "ПАР-СМЕСЬ-МАССИВ")
	}
}

// class composite

func (p *Parser) parseClassComposite(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "Композит класса"))
	}

	var n = &ast.ClassCompositeExpr{
		ExprBase: ast.ExprBase{Pos: p.pos},
		X:        x,
		Values:   make([]ast.ValuePair, 0),
	}

	p.expect(lexer.LBRACE)

	for p.tok != lexer.RBRACE && p.tok != lexer.EOF {

		var vp = ast.ValuePair{Pos: p.pos}

		vp.Name = p.parseIdent()
		p.expect(lexer.COLON)
		vp.Value = p.parseExpression()

		n.Values = append(n.Values, vp)

		if p.tok == lexer.RBRACE {
			break
		}
		p.expect(lexer.COMMA)
	}

	p.expect(lexer.RBRACE)

	return n
}

package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/lexer"
)

var _ = fmt.Printf

var validSimpleStmToken = tokens{
	lexer.IDENT: true,
	lexer.LPAR:  true,

	// literals
	lexer.INT:    true,
	lexer.FLOAT:  true,
	lexer.STRING: true,
	lexer.SYMBOL: true,

	// unary ops
	lexer.ADD: true,
	lexer.SUB: true,
	lexer.NOT: true,
}

var skipToStatement = tokens{
	lexer.EOF: true,

	lexer.RBRACE: true,

	lexer.IF:     true,
	lexer.WHILE:  true,
	lexer.RETURN: true,
	//TODO
}

var endStatementSeq = tokens{
	lexer.EOF:    true,
	lexer.RBRACE: true,
}

var endWhenCase = tokens{
	lexer.EOF:    true,
	lexer.RBRACE: true,

	lexer.IS:   true,
	lexer.ELSE: true,
}

//=== statements

func (p *Parser) parseStatementSeq() *ast.StatementSeq {
	if p.trace {
		defer un(trace(p, "Список операторов"))
	}

	if p.tok != lexer.LBRACE {
		p.expect(lexer.LBRACE)
		return &ast.StatementSeq{
			StatementBase: ast.StatementBase{Pos: p.pos},
			Statements:    make([]ast.Statement, 0),
		}
	}
	p.next()

	var n = p.parseStatementList(endStatementSeq)

	p.expect(lexer.RBRACE)

	return n
}

func (p *Parser) parseStatementList(stop tokens) *ast.StatementSeq {

	var n = &ast.StatementSeq{
		StatementBase: ast.StatementBase{Pos: p.pos},
		Statements:    make([]ast.Statement, 0),
	}

	for !stop[p.tok] {
		var s = p.parseStatement()
		if s != nil {
			n.Statements = append(n.Statements, s)
		}

		if stop[p.tok] {
			break
		}
		p.sep()
	}

	return n
}

func (p *Parser) parseStatement() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор"))
	}

	switch p.tok {
	case lexer.VAR:
		return &ast.DeclStatement{
			StatementBase: ast.StatementBase{Pos: p.pos},
			D:             p.parseVarDecl(),
		}
	case lexer.IF:
		return p.parseIf()
	case lexer.WHILE:
		return p.parseWhile()
	case lexer.WHEN:
		return p.parseWhen()
	case lexer.RETURN:
		return p.parseReturn()
	case lexer.BREAK:
		return p.parseBreak()
	case lexer.CRASH:
		return p.parseCrash()
	case lexer.GUARD:
		return p.parseGuard()

	default:
		if validSimpleStmToken[p.tok] {
			return p.parseSimpleStatement()
		}
		p.error(p.pos, "ПАР-ОШ-ОПЕРАТОР", p.tok.String())
		p.skipTo((skipToStatement))
	}

	return nil
}

func (p *Parser) parseSimpleStatement() ast.Statement {
	if p.trace {
		defer un(trace(p, "Простой оператор"))
	}

	var expr = p.parseExpression()

	switch p.tok {
	case lexer.ASSIGN:
		return p.parseAssign(expr)
	case lexer.INC:
		var n = &ast.IncStatement{
			StatementBase: ast.StatementBase{Pos: p.pos},
			L:             expr,
		}
		p.next()
		return n
	case lexer.DEC:
		var n = &ast.DecStatement{
			StatementBase: ast.StatementBase{Pos: p.pos},
			L:             expr,
		}
		p.next()
		return n

	default:
		var s = &ast.ExprStatement{
			StatementBase: ast.StatementBase{Pos: p.pos},

			X: expr,
		}

		return s
	}
}

func (p *Parser) parseAssign(l ast.Expr) ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор присваивания"))
	}

	var n = &ast.AssignStatement{
		StatementBase: ast.StatementBase{Pos: p.pos},
		L:             l,
	}

	p.next()
	n.R = p.parseExpression()

	return n
}

//====

func (p *Parser) parseIf() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор если"))
	}

	var n = &ast.If{
		StatementBase: ast.StatementBase{Pos: p.pos},
	}

	p.next()
	n.Cond = p.parseExpression()
	n.Then = p.parseStatementSeq()

	if p.tok != lexer.ELSE {
		return n
	}

	p.next()

	if p.tok == lexer.IF {
		n.Else = p.parseIf()
	} else {
		n.Else = p.parseStatementSeq()
	}

	return n
}

func (p *Parser) parseWhile() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор пока"))
	}

	var n = &ast.While{
		StatementBase: ast.StatementBase{Pos: p.pos},
	}

	p.next()
	n.Cond = p.parseExpression()
	n.Seq = p.parseStatementSeq()

	return n
}

func (p *Parser) parseReturn() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор вернуть"))
	}

	var n = &ast.Return{
		StatementBase: ast.StatementBase{Pos: p.pos},
	}

	p.next()

	if p.afterNL || p.tok == lexer.SEMI || p.tok == lexer.RBRACE {
		return n
	}

	n.X = p.parseExpression()

	return n
}

func (p *Parser) parseBreak() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор прервать"))
	}

	var n = &ast.Break{
		StatementBase: ast.StatementBase{Pos: p.pos},
	}

	p.next()

	return n
}

func (p *Parser) parseCrash() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор авария"))
	}

	var n = &ast.Crash{
		StatementBase: ast.StatementBase{Pos: p.pos},
	}

	p.next()
	p.expect(lexer.LPAR)
	n.X = p.parseExpression()
	p.expect(lexer.RPAR)

	return n
}

func (p *Parser) parseGuard() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор надо"))
	}

	var n = &ast.Guard{
		StatementBase: ast.StatementBase{Pos: p.pos},
	}

	p.next()

	n.Cond = p.parseExpression()

	p.expect(lexer.ELSE)

	switch p.tok {
	case lexer.RETURN:
		n.Else = p.parseReturn()
	case lexer.BREAK:
		n.Else = p.parseBreak()
	case lexer.CRASH:
		n.Else = p.parseCrash()
	default:
		n.Else = p.parseStatementSeq()
	}

	return n
}

func (p *Parser) parseWhen() ast.Statement {
	if p.trace {
		defer un(trace(p, "Оператор когда"))
	}

	var n = &ast.When{
		StatementBase: ast.StatementBase{Pos: p.pos},
	}

	p.next()
	n.X = p.parseExpression()
	p.expect(lexer.LBRACE)

	for p.tok == lexer.IS {
		var c = p.parseWhenCase()
		n.Cases = append(n.Cases, c)
	}

	if p.tok == lexer.ELSE {
		p.next()
		n.Else = p.parseStatementList(endStatementSeq)
	}
	p.expect(lexer.RBRACE)

	return n
}

func (p *Parser) parseWhenCase() *ast.Case {
	if p.trace {
		defer un(trace(p, "Оператор когда есть"))
	}

	var c = &ast.Case{
		StatementBase: ast.StatementBase{Pos: p.pos},
		Exprs:         make([]ast.Expr, 0),
	}
	p.next()

	for {
		var x = p.parseExpression()
		c.Exprs = append(c.Exprs, x)
		if p.tok != lexer.COMMA {
			break
		}
		p.next()
	}
	p.expect(lexer.COLON)

	c.Seq = p.parseStatementList(endWhenCase)

	return c
}

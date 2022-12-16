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
		p.sep()
	}

	p.expect(lexer.RBRACE)

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
	case lexer.CONST:
		panic("ni")

	case lexer.IF:
		return p.parseIf()
	case lexer.WHILE:
		return p.parseWhile()

	case lexer.RETURN:
		return p.parseReturn()

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

package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

//=== statements

func (p *Parser) parseFn() *ast.Function {
	if p.trace {
		defer un(trace(p, "Функция"))
	}

	var mod = ""
	if p.tok == lexer.MODIFIER {
		mod = p.lit
		p.next()
	}

	var n = &ast.Function{
		DeclBase: ast.DeclBase{Pos: p.pos},
	}

	p.expect(lexer.FN)

	n.Name = p.parseIdent()
	if p.parseExportMark() {
		n.SetExported()
	}

	n.Typ = p.parseFuncType()
	n.Seq = p.parseStatementSeq()

	switch mod {
	case "":
	case "@внешняя":
		n.External = true
	default:
		env.AddError(p.pos, "ПАР-ОШ-МОДИФИКАТОР", mod)
	}

	return n
}

func (p *Parser) parseFuncType() *ast.FuncType {
	if p.trace {
		defer un(trace(p, "Тип функции"))
	}

	var t = &ast.FuncType{
		TypeBase: ast.TypeBase{Pos: p.pos},
	}

	p.expect(lexer.LPAR)

	p.expect(lexer.RPAR)

	return t
}

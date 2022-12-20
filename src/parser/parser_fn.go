package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/lexer"
)

var _ = fmt.Printf

//=== statements

func (p *Parser) parseFn() *ast.Function {
	if p.trace {
		defer un(trace(p, "Функция"))
	}

	var n = &ast.Function{
		DeclBase: ast.DeclBase{Pos: p.pos},
	}

	p.expect(lexer.FN)

	if p.tok == lexer.LPAR { //receiver
		p.next()

		n.Recv = &ast.Param{
			TypeBase: ast.TypeBase{Pos: p.pos},
		}

		n.Recv.Name = p.parseIdent()
		p.expect(lexer.COLON)
		n.Recv.Typ = p.parseTypeRef()

		p.expect(lexer.RPAR)
	}

	n.Name = p.parseIdent()
	if p.parseExportMark() {
		n.SetExported()
	}

	n.Typ = p.parseFuncType()

	if p.tok == lexer.MODIFIER {
		var mod = p.lit
		p.next()

		switch mod {
		case "@внешняя":
			n.External = true
		default:
			p.error(p.pos, "ПАР-ОШ-МОДИФИКАТОР", mod)
		}

	} else {
		n.Seq = p.parseStatementSeq()
	}

	return n
}

func (p *Parser) parseFuncType() *ast.FuncType {
	if p.trace {
		defer un(trace(p, "Тип функции"))
	}

	var ft = &ast.FuncType{
		TypeBase: ast.TypeBase{Pos: p.pos},
	}

	p.expect(lexer.LPAR)

	p.parseParameters(ft)

	p.expect(lexer.RPAR)

	if p.tok == lexer.COLON {
		p.next()
		ft.ReturnTyp = p.parseTypeRef()
	}

	return ft
}

var skipToParam = tokens{
	lexer.EOF: true,

	lexer.RPAR:  true,
	lexer.COMMA: true,
}

func (p *Parser) parseParameters(ft *ast.FuncType) {

	for p.tok != lexer.RPAR && p.tok != lexer.EOF {

		var param = &ast.Param{
			TypeBase: ast.TypeBase{Pos: p.pos},
		}

		param.Name = p.parseIdent()

		p.expect(lexer.COLON)

		param.Typ = p.parseTypeRef()

		ft.Params = append(ft.Params, param)

		if p.tok == lexer.RPAR {
			break
		}
		p.expect(lexer.COMMA)
	}
}

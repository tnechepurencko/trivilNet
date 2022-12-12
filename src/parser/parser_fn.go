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

	var mod = ""
	if p.tok == lexer.MODIFIER {
		mod = p.lit
		p.next()
	}

	var n = &ast.Function{
		DeclBase: ast.DeclBase{Pos: p.pos},
	}

	switch mod {
	case "":
	case "@внешняя":
		n.External = true
	default:
		p.error(p.pos, "ПАР-ОШ-МОДИФИКАТОР", mod)
	}

	p.expect(lexer.FN)

	//receiver

	n.Name = p.parseIdent()
	if p.parseExportMark() {
		n.SetExported()
	}

	n.Typ = p.parseFuncType()

	if !n.External {
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
	}
}

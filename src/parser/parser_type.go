package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/lexer"
)

var _ = fmt.Printf

//=== types

func (p *Parser) parseTypeRef() ast.Type {
	if p.trace {
		defer un(trace(p, "Cсылка на тип"))
	}

	//TODO: мб
	var t = &ast.TypeRef{
		TypeBase: ast.TypeBase{Pos: p.pos},
	}

	var s = p.parseIdent()

	if p.tok == lexer.DOT {
		t.ModuleName = s
		t.TypeName = p.parseIdent()
	} else {
		t.TypeName = s
	}

	return t
}

func (p *Parser) parseTypeDecl() *ast.TypeDecl {
	if p.trace {
		defer un(trace(p, "Описание типа"))
	}

	p.next()

	var n = &ast.TypeDecl{
		DeclBase: ast.DeclBase{Pos: p.pos},
	}

	n.Name = p.parseIdent()
	if p.parseExportMark() {
		n.SetExported()
	}

	p.expect(lexer.EQ)
	n.Typ = p.parseTypeDef()

	return n
}

func (p *Parser) parseTypeDef() ast.Type {

	switch p.tok {
	case lexer.LBRACK:
		return p.parseArrayType()
	case lexer.CLASS:
		panic("ni")
	default:
		p.error(p.pos, "ПАР-ОШ-ОП-ТИПА", p.tok.String())
		return &ast.InvalidType{
			TypeBase: ast.TypeBase{Pos: p.pos},
		}
	}
}

func (p *Parser) parseArrayType() *ast.ArrayType {

	var t = &ast.ArrayType{
		TypeBase: ast.TypeBase{Pos: p.pos},
	}

	p.next()
	p.expect(lexer.RBRACK)

	t.ElementTyp = p.parseTypeRef()

	return t
}

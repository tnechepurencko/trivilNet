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

	if p.tok != lexer.IDENT {
		p.expect(lexer.IDENT)
	} else {
		var s = p.parseIdent()

		if p.tok == lexer.DOT {
			t.ModuleName = s
			t.TypeName = p.parseIdent()
		} else {
			t.TypeName = s
		}
	}

	return t
}

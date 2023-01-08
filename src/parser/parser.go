package parser

import (
	"fmt"
	"strconv"

	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

const ExportMark = lexer.MUL

type Parser struct {
	source *env.Source
	lex    *lexer.Lexer
	module *ast.Module

	pos     int
	tok     lexer.Token
	lit     string
	afterNL bool

	trace  bool
	indent int
}

func Parse(source *env.Source) *ast.Module {
	p := new(Parser)
	p.lex = new(lexer.Lexer)
	p.lex.Init(source)
	p.source = source

	p.trace = *env.TraceParser

	p.next()
	p.parseModule()

	return p.module
}

func (p *Parser) next() {
	p.afterNL = false
	for true {
		p.pos, p.tok, p.lit = p.lex.Scan()

		switch p.tok {
		case lexer.EOF:
			p.afterNL = true
			return
		case lexer.BLOCK_COMMENT:
			continue
		case lexer.LINE_COMMENT:
			p.afterNL = true
			continue
		case lexer.NL:
			p.afterNL = true
			continue
		default:
			return
		}
	}
}

//====

func (p *Parser) error(pos int, id string, args ...interface{}) {
	s := env.AddError(pos, id, args...)
	if p.trace {
		fmt.Println(s)
	}
}

func (p *Parser) expect(tok lexer.Token) {
	if p.tok != tok {
		p.error(p.pos, "ПАР-ОЖИДАЛСЯ", tok.String())
	}
	p.next()
}

func (p *Parser) sep() {
	if p.tok == lexer.SEMI {
		p.next()
	} else if p.afterNL {
		// ok
	} else {
		p.error(p.pos, "ПАР-ОШ-РАЗД", p.tok.String())
	}
}

//====

type tokens map[lexer.Token]bool

func (p *Parser) skipTo(ts tokens) {
	p.next()

	for {
		if _, ok := ts[p.tok]; ok {
			break
		}
		p.next()
	}
}

//=====

func (p *Parser) parseModule() {

	p.module = ast.NewModule()

	if p.tok != lexer.MODULE {
		p.error(p.pos, "ПАР-ОЖИДАЛСЯ", lexer.MODULE.String())
		return
	}
	p.next()
	p.module.Name = p.parseIdent()
	p.sep()

	p.parseImportList()
	p.parseDeclarations()
}

func (p *Parser) parseImportList() {
	if p.trace {
		defer un(trace(p, "Импорты"))
	}

	for p.tok == lexer.IMPORT {

		var n = &ast.Import{Pos: p.pos}

		p.next()
		if p.tok == lexer.STRING {
			n.Path = p.lit
			p.next()
		} else {
			p.expect(lexer.STRING)
		}

		p.module.Imports = append(p.module.Imports, n)

		p.sep()
	}
}

//====

var skipToDeclaration = tokens{
	lexer.EOF: true,

	lexer.TYPE:  true,
	lexer.VAR:   true,
	lexer.CONST: true,
	lexer.FN:    true,
	lexer.ENTRY: true,
}

func (p *Parser) parseDeclarations() {
	if p.trace {
		defer un(trace(p, "Описания"))
	}

	for p.tok != lexer.EOF {

		switch p.tok {
		case lexer.TYPE:
			var d = p.parseTypeDecl()
			d.Host = p.module
			p.module.Decls = append(p.module.Decls, d)
		case lexer.CONST:
			var cs = p.parseConstDecls()
			for _, c := range cs {
				c.Host = p.module
				p.module.Decls = append(p.module.Decls, c)
			}
		case lexer.VAR:
			var d = p.parseVarDecl()
			d.Host = p.module
			p.module.Decls = append(p.module.Decls, d)
		case lexer.FN:
			var d = p.parseFn()
			d.Host = p.module
			p.module.Decls = append(p.module.Decls, d)
		case lexer.ENTRY:
			p.parseEntry()
		default:
			p.error(p.pos, "ПАР-ОШ-ОПИСАНИЕ", p.tok.String())
			p.skipTo(skipToDeclaration)
			continue
		}

		p.sep()
	}
}

//=== константы

func (p *Parser) parseConstDecls() []*ast.ConstDecl {
	if p.trace {
		defer un(trace(p, "Описание констант"))
	}

	p.next()

	if p.tok == ExportMark || p.tok == lexer.LPAR {
		return p.parseConstGroup()
	}

	var c = p.parseSingleConst()

	var cs = make([]*ast.ConstDecl, 1)
	cs[0] = c

	return cs
}

func (p *Parser) parseSingleConst() *ast.ConstDecl {

	var n = &ast.ConstDecl{
		DeclBase: ast.DeclBase{Pos: p.pos},
	}

	n.Name = p.parseIdent()
	if p.parseExportMark() {
		n.Exported = true
	}

	p.expect(lexer.COLON)
	n.Typ = p.parseTypeRef()

	p.expect(lexer.EQ)
	n.Value = p.parseExpression() //! const expression

	return n
}

func (p *Parser) parseNextConst() *ast.ConstDecl {

	var n = &ast.ConstDecl{
		DeclBase: ast.DeclBase{Pos: p.pos},
	}

	n.Name = p.parseIdent()
	if p.parseExportMark() {
		n.Exported = true
	}

	if p.tok != lexer.COLON {
		return n
	}

	p.next()
	n.Typ = p.parseTypeRef()

	p.expect(lexer.EQ)
	n.Value = p.parseExpression()

	return n
}

func (p *Parser) parseConstGroup() []*ast.ConstDecl {

	var exported = false
	if p.tok == ExportMark {
		exported = true
		p.next()
	}

	p.expect(lexer.LPAR)

	var cs = make([]*ast.ConstDecl, 0)
	var c = p.parseSingleConst()

	c.Exported = exported
	cs = append(cs, c)

	for p.tok != lexer.RPAR && p.tok != lexer.EOF {
		p.sep()

		c = p.parseNextConst()
		cs = append(cs, c)
	}

	p.expect(lexer.RPAR)

	p.completeConstGroup(cs)

	return cs
}

func (p *Parser) completeConstGroup(cs []*ast.ConstDecl) {

	var base *ast.ConstDecl
	var first = true
	var val int

	for _, c := range cs {
		if c.Typ != nil {
			base = c
			first = true
		} else {
			if first { // первая авто константа

				literal, ok := base.Value.(*ast.LiteralExpr)

				if ok && literal.Kind == ast.Lit_Int {
					i, err := strconv.Atoi(literal.Lit)
					if err != nil {
						p.error(base.Pos, "ПАР-ОШ-КОНСТ-БАЗА", fmt.Sprintf("(%s)", err.Error()))
						val = 0
					} else {
						val = i
					}
				} else {
					p.error(base.Pos, "ПАР-ОШ-КОНСТ-БАЗА", "")
					val = 0
				}
				first = false
			}

			c.Typ = base.Typ

			val++
			c.Value = &ast.LiteralExpr{
				ExprBase: ast.ExprBase{Pos: c.Pos},
				Kind:     ast.Lit_Int,
				Lit:      strconv.Itoa(val),
			}
		}
	}

}

//=== переменные

func (p *Parser) parseVarDecl() *ast.VarDecl {
	if p.trace {
		defer un(trace(p, "Описание переменной"))
	}

	p.next()

	var n = &ast.VarDecl{
		DeclBase: ast.DeclBase{Pos: p.pos},
	}

	n.Name = p.parseIdent()
	if p.parseExportMark() {
		n.Exported = true
	}

	if p.tok == lexer.COLON {
		p.next()
		n.Typ = p.parseTypeRef()
	}

	if p.tok != lexer.EQ {
		env.AddError(p.pos, "ПАР-ПЕРЕМ-ИНИТ")
		return n
	}

	p.expect(lexer.EQ)
	n.Init = p.parseExpression()

	return n
}

//==== вход

func (p *Parser) parseEntry() {
	if p.trace {
		defer un(trace(p, "Вход"))
	}

	var n = &ast.EntryFn{
		Pos: p.pos,
	}

	p.next()
	n.Seq = p.parseStatementSeq()

	if p.module.Entry != nil {
		p.error(p.pos, "ПАР-ДУБЛЬ-ВХОД")
		return
	}

	p.module.Entry = n
}

//====

func (p *Parser) parseIdent() string {
	name := "_"
	if p.tok == lexer.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(lexer.IDENT)
	}
	return name
}

func (p *Parser) parseExportMark() bool {
	if p.tok == ExportMark {
		p.next()
		return true
	}
	return false
}

//====

type modifier struct {
	name   string
	attrs  []string
	values []string
}

func (p *Parser) parseModifier() modifier {

	var mod = modifier{
		name: p.lit,
	}
	p.next()

	if p.tok != lexer.LPAR {
		return mod
	}

	p.next()

	for p.tok != lexer.RPAR && p.tok != lexer.EOF {
		var attr = ""
		var value = ""

		if p.tok == lexer.STRING {
			attr = p.lit
			p.next()
		} else {
			p.expect(lexer.STRING)
		}

		p.expect(lexer.COLON)

		if p.tok == lexer.STRING {
			value = p.lit
			p.next()
		} else {
			p.expect(lexer.STRING)
		}

		mod.attrs = append(mod.attrs, attr)
		mod.values = append(mod.values, value)

		if p.tok == lexer.RPAR {
			break
		}
		p.expect(lexer.COMMA)
	}
	p.expect(lexer.RPAR)

	return mod
}

package parser

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

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

	p.trace = *env.TraceFlag

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

func (p *Parser) expectSep(err string) {
	if p.tok == lexer.SEMI {
		p.next()
	} else if p.afterNL {
		// ok
	} else {
		p.error(p.pos, err, p.tok.String())
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

	p.parseImportList()
	p.parseDeclarations()
}

func (p *Parser) parseImportList() {
	if p.trace {
		defer un(trace(p, "Импорты"))
	}

}

//====

var skipToDeclaration = tokens{
	lexer.EOF: true,

	lexer.TYPE:     true,
	lexer.VAR:      true,
	lexer.CONST:    true,
	lexer.FN:       true,
	lexer.MODIFIER: true,
	lexer.ENTRY:    true,
}

func (p *Parser) parseDeclarations() {
	if p.trace {
		defer un(trace(p, "Описания"))
	}

	var d ast.Decl

	for p.tok != lexer.EOF {

		d = nil
		switch p.tok {
		case lexer.FN, lexer.MODIFIER:
			d = p.parseFn()
		case lexer.ENTRY:
			p.parseEntry()
		default:
			p.error(p.pos, "ПАР-ОШ-ОПИСАНИЕ", p.tok.String())
			p.skipTo(skipToDeclaration)
			continue
		}

		p.expectSep("ПАР_РАЗД_ОПИСАНИЙ")

		if d != nil {
			p.module.Decls = append(p.module.Decls, d)
		}

	}
}

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
	if p.tok == lexer.MUL {
		p.next()
		return true
	}
	return false
}

package lexer

import (
	"strconv"
	//"unicode"
)

type Token int

const (
	// Special tokens
	Invalid Token = iota
	EOF
	NL

	LINE_COMMENT
	BLOCK_COMMENT
	MODIFIER

	// literals
	IDENT
	INT
	FLOAT
	STRING

	// operators
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND // &
	OR  // |
	NOT // ~

	BITAND // &.
	BITOR  // |.
	INC    // ++
	DEC    // --

	EQ  // =
	LSS // <
	GTR // >

	NEQ // #
	LEQ // <=
	GEQ // >=

	NNQUERY // ?
	NNCHECK // !

	ELLIPSIS // ...
	ASSIGN   // :=

	LPAR   // (
	RPAR   // )
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	RBRACE // }
	LCONV  // (:
	//RCONV  // ›

	COMMA // ,
	DOT   // .
	SEMI  // ;
	COLON // :

	// keywords
	keyword_beg
	BREAK
	CONST
	CRASH
	ELSE
	ENTRY
	FN
	GUARD
	IF
	IMPORT
	MODULE
	RETURN
	CLASS
	TYPE
	VAR
	WHILE

	keyword_end
)

var tokens = [...]string{
	Invalid: "Invalid",

	EOF: "EOF",
	NL:  "NL",

	LINE_COMMENT:  "LINE_COMMENT",
	BLOCK_COMMENT: "BLOCK_COMMENT",

	MODIFIER: "@",

	IDENT:  "идентификатор",
	INT:    "целый литерал",
	FLOAT:  "вещественный литерал",
	STRING: "строковый литерал",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	NOT: "~",
	//	XOR:     "^",

	BITAND: "&.",
	BITOR:  "|.",

	INC: "++",
	DEC: "--",

	EQ:  "=",
	LSS: "<",
	GTR: ">",

	NEQ: "#",
	LEQ: "<=",
	GEQ: ">=",

	NNQUERY: "?", // так ли?
	NNCHECK: "!", // так ли?

	ELLIPSIS: "...",
	ASSIGN:   ":=",

	LPAR: "(",
	RPAR: ")",

	LBRACK: "[",
	RBRACK: "]",

	LBRACE: "{",
	RBRACE: "}",

	LCONV: "(:", // or .<
	//RCONV: "›",

	COMMA: ",",
	DOT:   ".",
	SEMI:  ";",
	COLON: ":",

	CRASH:  "авария",
	RETURN: "вернуть",
	ENTRY:  "вход",
	IF:     "если",
	ELSE:   "иначе",
	IMPORT: "импорт",
	CLASS:  "класс",
	CONST:  "конст",
	MODULE: "модуль",
	GUARD:  "надо",
	WHILE:  "пока",
	BREAK:  "прервать",
	VAR:    "пусть",
	TYPE:   "тип",
	FN:     "фн",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_beg+1))
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

/*
func IsKeyword(name string) bool {
	// TODO: opt: use a perfect hash function instead of a global map.
	_, ok := keywords[name]
	return ok
}

func IsIdentifier(name string) bool {
	if name == "" || IsKeyword(name) {
		return false
	}
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}
*/

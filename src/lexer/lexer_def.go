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

	ASSIGN // :=

	LPAR   // (
	RPAR   // )
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	RBRACE // }

	COMMA  // ,
	PERIOD // .
	SEMI   // ;
	COLON  // :

	// keywords
	keyword_beg
	BREAK
	CONST
	ELSE
	FN
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

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

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

	ASSIGN: ":=",

	LPAR:   "(",
	LBRACK: "[",
	LBRACE: "{",

	RPAR:   ")",
	RBRACK: "]",
	RBRACE: "}",

	COMMA:  ",",
	PERIOD: ".",
	SEMI:   ";",
	COLON:  ":",

	BREAK:  "прервать",
	CLASS:  "класс",
	CONST:  "конст",
	ELSE:   "иначе",
	FN:     "фн",
	IF:     "если",
	IMPORT: "импорт",
	MODULE: "модуль",
	RETURN: "вернуть",
	TYPE:   "тип",
	VAR:    "пусть",
	WHILE:  "пока",
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

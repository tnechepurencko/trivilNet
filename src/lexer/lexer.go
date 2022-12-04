package lexer

import (
	//"bytes"
	"fmt"
	//"strconv"
	"unicode"
	"unicode/utf8"

	"trivil/env"
)

var _ = fmt.Printf

type Lexer struct {
	source *env.Source
	src    []byte

	// состояние
	ch         rune // текущий символ
	offset     int  // смещение текущего
	rdOffset   int  // позиция чтения
	lineOffset int  // смещение текущей строки

}

func (s *Lexer) Init(source *env.Source) {
	s.source = source
	s.src = source.Bytes

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0

	// пропустить BOM
	if len(s.src) >= 3 && s.src[0] == 0xEF && s.src[1] == 0xBB && s.src[2] == 0xBF {
		s.offset = 3
	}

	s.next()
}

func (s *Lexer) error(ofs int, id string, args ...interface{}) {
	env.AddError(s.source.MakePos(ofs), id, args...)
}

// Read the next Unicode char into s.ch.
// s.ch < 0 means end-of-file.
func (s *Lexer) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.source.AddLine(s.offset)
		}
		r := rune(s.src[s.rdOffset])
		w := 1
		switch {
		case r == 0:
			s.error(s.offset, "ЛЕК-ОШ-СИМ", rune(0))
		case r == '\r':
			r = '\n'
			if s.rdOffset+1 < len(s.src) && rune(s.src[s.rdOffset+1]) == '\n' {
				w = 2
			}
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "ЛЕК-UTF8")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.source.AddLine(s.offset)
		}
		s.ch = -1 // eof
	}
}

/* может пригодится
func (s *Lexer) peek() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}
	return 0
}
*/

//====

func lower(ch rune) rune     { return ('a' - 'A') | ch }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func (s *Lexer) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' {
		s.next()
	}
}

func (s *Lexer) scanModifier() string {
	ofs := s.offset - 1
	if !isLetter(s.ch) {
		s.error(s.offset, "ЛЕК-МОДИФИКАТОР")
		return ""
	}
	for isLetter(s.ch) {
		s.next()
	}
	return string(s.src[ofs:s.offset])
}

func (s *Lexer) scanLineComment() int {
	// Первый '/' уже взят
	ofs := s.offset - 1
	s.next()
	for s.ch != '\n' && s.ch >= 0 {
		s.next()
	}
	if s.ch == '\n' {
		s.next() // перешли на след. символ после комментария
	}
	return ofs
}

func (s *Lexer) scanBlockComment() int {
	// '/' уже взят
	ofs := s.offset - 1
	s.next() // '*'
	for true {
		if s.ch < 0 {
			s.error(ofs, "ЛЕК-НЕТ-*/")
			break
		}
		ch := s.ch
		s.next()
		if ch == '*' && s.ch == '/' {
			s.next()
			break
		} else if ch == '/' && s.ch == '*' {
			s.scanBlockComment()
		}
	}
	return ofs
}

// пока упрощенный
func (s *Lexer) scanIdentifier() string {
	ofs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[ofs:s.offset])
}

func (s *Lexer) scanString(opening rune) string {
	// первая кавычка уже взята
	ofs := s.offset - 1

	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			s.error(ofs, "ЛЕК-ОШ-СТРОКА")
			break
		}
		s.next()
		if ch == opening {
			break
		}
		if ch == '\\' {
			s.scanEscape(opening)
		}
	}

	return string(s.src[ofs:s.offset])
}

// Сканирует escape sequence. В случае ошибки возвращает false
// Не проверяет корректность
func (s *Lexer) scanEscape(quote rune) bool {
	ofs := s.offset

	var n int
	if s.ch == 'u' { // \uABCD
		n = 4
	} else {
		if s.ch < 0 {
			s.error(ofs, "ЛЕК-ОШ-ESCAPE")
			return false
		}
		s.next()
		return true

	}

	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= 16 {
			if s.ch < 0 {
				s.error(s.offset, "ЛЕК-ОШ-ESCAPE")
			} else {
				s.error(s.offset, "ЛЕК-ОШ-СИМ", s.ch)
				return false
			}
		}
		s.next()
		n--
	}

	return true
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16 // larger than any legal digit val
}

//==== number

func (s *Lexer) scanDigits(base int) {
	if base == 10 {
		for isDecimal(s.ch) {
			s.next()
		}
	} else if base == 16 {
		for isHex(s.ch) {
			s.next()
		}
	} else {
		panic("! wrong base")
	}
}

func (s *Lexer) scanNumber() (Token, string) {
	ofs := s.offset

	base := 10
	prefix := rune(0) // 0 - нет префикса, 'x' - 0x

	// целое
	if s.ch == '0' {
		s.next()
		if s.ch == 'x' {
			s.next()
			base = 16
			prefix = 'x'
		}

	}

	s.scanDigits(base)

	if s.ch != '.' {
		return INT, string(s.src[ofs:s.offset])
	}

	// дробная часть
	if prefix != rune(0) {
		s.error(s.offset, "ЛЕК-ВЕЩ-БАЗА")
	}

	s.next()
	s.scanDigits(10)

	return FLOAT, string(s.src[ofs:s.offset])
}

//====

func (s *Lexer) checkEqu(tok0, tok1 Token) Token {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	return tok0
}

func (s *Lexer) checkNext(tok0 Token, next rune, tok1 Token) Token {
	if s.ch == next {
		s.next()
		return tok1
	}
	return tok0
}

func (s *Lexer) Scan() (pos int, tok Token, lit string) {

	s.skipWhitespace()

	// начало лексемы
	pos = s.source.MakePos(s.offset)

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		if len(lit) > 1 {
			// все ключевые слова длиннее одного символа
			tok = Lookup(lit)
		} else {
			tok = IDENT
		}
	case isDecimal(ch):
		tok, lit = s.scanNumber()
	default:
		s.next() // всегда двигаемся
		switch ch {
		case -1:
			tok = EOF
		case '\n':
			// пропускаем все
			for s.ch == '\n' {
				s.next()
			}
			tok = NL
		case '"':
			tok = STRING
			lit = s.scanString('"')
		case '\'':
			tok = STRING
			lit = s.scanString('\'')
		case '@':
			tok = MODIFIER
			lit = s.scanModifier()
		case ':':
			tok = s.checkEqu(COLON, ASSIGN)
		case '.':
			tok = DOT
		case ',':
			tok = COMMA
		case ';':
			tok = SEMI
			lit = ";"
		case '(':
			tok = LPAR
		case ')':
			tok = RPAR
		case '[':
			tok = LBRACK
		case ']':
			tok = RBRACK
		case '{':
			tok = LBRACE
		case '}':
			tok = RBRACE
		case '+':
			tok = s.checkNext(ADD, '+', INC)
		case '-':
			tok = s.checkNext(SUB, '-', DEC)
		case '*':
			tok = MUL
		case '/':
			if s.ch == '/' {
				tok = LINE_COMMENT
				ofs := s.scanLineComment()
				lit = string(s.src[ofs:s.offset])
			} else if s.ch == '*' {
				tok = BLOCK_COMMENT
				ofs := s.scanBlockComment()
				lit = string(s.src[ofs:s.offset])
			} else {
				tok = QUO
			}

		case '%':
			tok = REM
			//		case '^':
			//			tok = XOR
		case '<':
			tok = s.checkEqu(LSS, LEQ)
		case '>':
			tok = s.checkEqu(GTR, GEQ)
		case '=':
			tok = EQ
		case '#':
			tok = NEQ
		case '~':
			tok = NOT
		case '&':
			tok = s.checkNext(AND, '.', BITAND)
		case '|':
			tok = s.checkNext(OR, '.', BITOR)
		default:
			s.error(s.offset, "ЛЕК-ОШ-СИМ", ch)
			tok = Invalid
			lit = string(ch)
		}
	}

	return
}

func (s *Lexer) WhitespaceBefore(c rune) bool {

	ofs := s.offset - 2
	if ofs < 0 {
		return false
	}
	if s.src[ofs+1] != byte(c) {
		return false
	}

	ch := s.src[ofs]
	return ch == ' ' || ch == '\t'
}

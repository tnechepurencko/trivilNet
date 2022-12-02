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

func (s *Lexer) error(pos int, id string, args ...interface{}) {
	env.AddError(s.source, pos, id, args...)
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

func (s *Lexer) peek() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}
	return 0
}

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

/*

func trailingDigits(text []byte) (int, int, bool) {
	i := bytes.LastIndexByte(text, ':') // look from right (Windows filenames may contain ':')
	if i < 0 {
		return 0, 0, false // no ":"
	}
	// i >= 0
	n, err := strconv.ParseUint(string(text[i+1:]), 10, 0)
	return i + 1, int(n), err == nil
}

func (s *Lexer) findLineEnd() bool {
	// initial '/' already consumed

	defer func(offs int) {
		// reset scanner state to where it was upon calling findLineEnd
		s.ch = '/'
		s.offset = offs
		s.rdOffset = offs + 1
		s.next() // consume initial '/' again
	}(s.offset - 1)

	// read ahead until a newline, EOF, or non-comment token is found
	for s.ch == '/' || s.ch == '*' {
		if s.ch == '/' {
			//-style comment always contains a newline
			return true
		}
		/*-style comment: look for newline *!/
		s.next()
		for s.ch >= 0 {
			ch := s.ch
			if ch == '\n' {
				return true
			}
			s.next()
			if ch == '*' && s.ch == '/' {
				s.next()
				break
			}
		}
		s.skipWhitespace() // s.insertSemi is set
		if s.ch < 0 || s.ch == '\n' {
			return true
		}
		if s.ch != '/' {
			// non-comment token
			return false
		}
		s.next() // consume '/'
	}

	return false
}




// digits accepts the sequence { digit | '_' }.
// If base <= 10, digits accepts any decimal digit but records
// the offset (relative to the source start) of a digit >= base
// in *invalid, if *invalid < 0.
// digits returns a bitset describing whether the sequence contained
// digits (bit 0 is set), or separators '_' (bit 1 is set).
func (s *Lexer) digits(base int, invalid *int) (digsep int) {
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			} else if s.ch >= max && *invalid < 0 {
				*invalid = s.offset // record invalid rune offset
			}
			digsep |= ds
			s.next()
		}
	} else {
		for isHex(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			}
			digsep |= ds
			s.next()
		}
	}
	return
}

func (s *Lexer) scanNumber() (Token, string) {
	offs := s.offset
	tok := ILLEGAL

	base := 10        // number base
	prefix := rune(0) // one of 0 (decimal), '0' (0-octal), 'x', 'o', or 'b'
	digsep := 0       // bit 0: digit present, bit 1: '_' present
	invalid := -1     // index of invalid digit in literal, or < 0

	// integer part
	if s.ch != '.' {
		tok = INT
		if s.ch == '0' {
			s.next()
			switch lower(s.ch) {
			case 'x':
				s.next()
				base, prefix = 16, 'x'
			case 'o':
				s.next()
				base, prefix = 8, 'o'
			case 'b':
				s.next()
				base, prefix = 2, 'b'
			default:
				base, prefix = 8, '0'
				digsep = 1 // leading 0
			}
		}
		digsep |= s.digits(base, &invalid)
	}

	// fractional part
	if s.ch == '.' {
		tok = FLOAT
		if prefix == 'o' || prefix == 'b' {
			s.error(s.offset, "invalid radix point in "+litname(prefix))
		}
		s.next()
		digsep |= s.digits(base, &invalid)
	}

	if digsep&1 == 0 {
		s.error(s.offset, litname(prefix)+" has no digits")
	}

	// exponent
	if e := lower(s.ch); e == 'e' || e == 'p' {
		switch {
		case e == 'e' && prefix != 0 && prefix != '0':
			s.errorf(s.offset, "%q exponent requires decimal mantissa", s.ch)
		case e == 'p' && prefix != 'x':
			s.errorf(s.offset, "%q exponent requires hexadecimal mantissa", s.ch)
		}
		s.next()
		tok = FLOAT
		if s.ch == '+' || s.ch == '-' {
			s.next()
		}
		ds := s.digits(10, nil)
		digsep |= ds
		if ds&1 == 0 {
			s.error(s.offset, "exponent has no digits")
		}
	} else if prefix == 'x' && tok == FLOAT {
		s.error(s.offset, "hexadecimal mantissa requires a 'p' exponent")
	}

	// suffix 'i'
	if s.ch == 'i' {
		tok = IMAG
		s.next()
	}

	lit := string(s.src[offs:s.offset])
	if tok == INT && invalid >= 0 {
		s.errorf(invalid, "invalid digit %q in %s", lit[invalid-offs], litname(prefix))
	}
	if digsep&2 != 0 {
		if i := invalidSep(lit); i >= 0 {
			s.error(offs+i, "'_' must separate successive digits")
		}
	}

	return tok, lit
}

func litname(prefix rune) string {
	switch prefix {
	case 'x':
		return "hexadecimal literal"
	case 'o', '0':
		return "octal literal"
	case 'b':
		return "binary literal"
	}
	return "decimal literal"
}

// invalidSep returns the index of the first invalid separator in x, or -1.
func invalidSep(x string) int {
	x1 := ' ' // prefix char, we only care if it's 'x'
	d := '.'  // digit, one of '_', '0' (a digit), or '.' (anything else)
	i := 0

	// a prefix counts as a digit
	if len(x) >= 2 && x[0] == '0' {
		x1 = lower(rune(x[1]))
		if x1 == 'x' || x1 == 'o' || x1 == 'b' {
			d = '0'
			i = 2
		}
	}

	// mantissa and exponent
	for ; i < len(x); i++ {
		p := d // previous digit
		d = rune(x[i])
		switch {
		case d == '_':
			if p != '0' {
				return i
			}
		case isDecimal(d) || x1 == 'x' && isHex(d):
			d = '0'
		default:
			if p == '_' {
				return i - 1
			}
			d = '.'
		}
	}
	if d == '_' {
		return len(x) - 1
	}

	return -1
}




func stripCR(b []byte, comment bool) []byte {
	c := make([]byte, len(b))
	i := 0
	for j, ch := range b {
		// In a /*-style comment, don't strip \r from *\r/ (incl.
		// sequences of \r from *\r\r...\r/) since the resulting
		// *!/ would terminate the comment too early unless the \r
		// is immediately following the opening /* in which case
		// it's ok because /*!/ is not closed yet (issue #11151).

		if ch != '\r' || comment && i > len("/*") && c[i-1] == '*' && j+1 < len(b) && b[j+1] == '/' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}

*/

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
	pos = s.offset

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		if len(lit) > 1 {
			// все ключевые слова длиннее одного символа
			tok = Lookup(lit)
		} else {
			tok = IDENT
		}
		//!	case isDecimal(ch) || ch == '.' && isDecimal(rune(s.peek())):
		//		tok, lit = s.scanNumber()
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

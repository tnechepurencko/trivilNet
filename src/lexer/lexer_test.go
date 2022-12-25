package lexer

import (
	"fmt"
	"testing"
	"trivil/env"
)

type pair struct {
	tok Token
	lit string
}

type one struct {
	text  string
	pairs []pair
}

var tests = []one{
	{"модуль", []pair{{MODULE, "модуль"}}},
	{"модуль ", []pair{{MODULE, "модуль"}}},
	{"модуль\n", []pair{{MODULE, "модуль"}, {NL, ""}}},
	{"сложное имя", []pair{{IDENT, "сложное имя"}}},
	{"№ сложное имя", []pair{{IDENT, "№ сложное имя"}}},
	{"фн сложное имя(", []pair{{FN, "фн"}, {IDENT, "сложное имя"}, {LPAR, ""}}},
	{"имя153", []pair{{IDENT, "имя153"}}},
	{"имя 153", []pair{{IDENT, "имя"}, {INT, "153"}}},
	{"как дела?", []pair{{IDENT, "как дела?"}}},
	{"как дела ?", []pair{{IDENT, "как дела"}, {NNQUERY, ""}}},
	{"Паниковать !", []pair{{IDENT, "Паниковать"}, {NNCHECK, ""}}},
	{"если-нет", []pair{{IDENT, "если-нет"}}},
	{"ц--", []pair{{IDENT, "ц"}, {DEC, ""}}},
}

//===

func TestValid(t *testing.T) {
	fmt.Printf("--- valid tests: %d ---\n", len(tests))
	t.Run("valid tests", func(t *testing.T) {
		for _, test := range tests {
			check(t, test)
		}
	})
}

func check(t *testing.T, test one) {

	var src = env.AddImmSource(test.text)
	var lex = new(Lexer)
	lex.Init(src)

	var actual = make([]pair, 0)
	for true {
		_, tok, lit := lex.Scan()
		if tok == EOF {
			break
		}
		actual = append(actual, pair{tok, lit})
	}

	for i := 0; i < len(actual) && i < len(test.pairs); i++ {
		if actual[i].tok != test.pairs[i].tok {
			t.Errorf("Лексема %s вместо %s в тексте\n'%s'", actual[i].tok, test.pairs[i].tok, test.text)
		} else if actual[i].lit != test.pairs[i].lit {
			t.Errorf("Текст лексемы '%s' вместо '%s' в тексте\n'%s'", actual[i].lit, test.pairs[i].lit, test.text)
		}
	}

	if len(actual) != len(test.pairs) {
		t.Errorf("Получено %d лексем вместо %d в тексте\n'%s'", len(actual), len(test.pairs), test.text)
	}
}

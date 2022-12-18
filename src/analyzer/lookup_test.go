package analyzer

import (
	"testing"
	"trivil/env"
	"trivil/parser"
)

var valid_texts = []string{
	"модуль м; пусть я: Цел; вход { я = 1 }",
	"модуль м; вход { пусть я: Цел; я = 1 }",
	"модуль м; конст ц: Цел = 1; вход { пусть я: Цел; я = ц }",
	"модуль м; фн Ф() {}; вход { Ф() }",
	"модуль м; тип Я = []Цел; пусть я: Я; вход { }",
	"модуль м; тип Я = []Цел; пусть я: Я; вход { я[0] := 1 }",
	"модуль м; тип Я = []Цел; вход { пусть я: Я; я[0] := 1 }",
	"модуль м; тип Я = []Цел; вход { пусть я: Я; я = Я[] }",
	"модуль м; тип Я = []Цел; вход { пусть я: Я; пусть й: Цел; я = Я[й] }",
	"модуль м; тип Я = []Цел; вход { пусть я: Я; пусть й: Цел; я = Я[й, й] }",
}

var invalid_texts = []string{
	"иначе м",
}

//===

func TestValid(t *testing.T) {
	t.Run("valid tests", func(t *testing.T) {
		for _, text := range valid_texts {
			checkValid(t, text)
		}
	})
}

func checkValid(t *testing.T, text string) {
	compile(text)
	if env.ErrorCount() > 0 {
		t.Errorf("Unexpected %d errors in text:\n%s\n%s\n", env.ErrorCount(), text, env.GetError(0))
		env.ClearErrors()
	}
}

/*
func TestInvalid(t *testing.T) {
	t.Run("invalid tests", func(t *testing.T) {
		for _, text := range invalid_texts {
			checkInvalid(t, text)
		}
	})
}

func checkInvalid(t *testing.T, text string) {
	parseSrc(text)
	if env.ErrorCount() > 0 {
		env.ClearErrors()
	} else {
		t.Errorf("Error(s) expected in text:\n%s\n", text)
	}
}
*/

func compile(text string) {
	var src = env.AddImmSource(text)

	m := parser.Parse(src)
	if env.ErrorCount() > 0 {
		return
	}
	lookup(m)
}

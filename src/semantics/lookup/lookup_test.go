package lookup

import (
	"fmt"
	"testing"
	"trivil/env"
	"trivil/parser"
)

var valid_texts = []string{
	"модуль м; пусть я: Цел64 = 0; вход { я = 1 }",
	"модуль м; вход { пусть я: Цел64 = 0; я = 1 }",
	"модуль м; конст ц: Цел64 = 1; вход { пусть я: Цел64  = 0; я = ц }",
	"модуль м; фн Ф() {}; вход { Ф() }",
	"модуль м; тип Я = []Цел64; пусть я: Я = 0; вход { }",
	"модуль м; тип Я = []Цел64; пусть я: Я = 0; вход { я[0] := 1 }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я = 0; я[0] := 1 }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я = 0; я = Я[] }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я = 0; пусть й: Цел64 = 0; я = Я[й] }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я = 0; пусть й: Цел64 = 0; я = Я[й, й] }",
	"модуль м; тип Я = []Цел64; тип М = []Я; вход { }",
	"модуль м; тип М = []Я; тип Я = []Цел64; вход { }",
	"модуль м; тип К = класс {}; вход { }",
	"модуль м; тип К = класс { я: Цел64 }; вход { }",
	"модуль м; тип К = класс { я: Цел64; к: К }; вход { }",
	"модуль м; тип К = класс { я: Цел64; к: К }; вход { К{} }",
	"модуль м; тип К = класс { я: Цел64; б: Цел64 }; вход { К{я: 1} }",
	"модуль м; тип К = класс { я: Цел64; б: Цел64 }; вход { К{я: 1, б: 2} }",
}

type one struct {
	text string
	id   string
}

var error_tests = []one{
	{"модуль м; вход { я := 1 }", "СЕМ-НЕ-НАЙДЕНО"},
	{"модуль м; вход { пусть я: Т := 1 }", "СЕМ-НЕ-НАЙДЕНО"},

	{"модуль м; вход { пусть я := 1; пусть я = 2 }", "СЕМ-УЖЕ-ОПИСАНО"},

	{"модуль м; вход { пусть я := 1; пусть ф: я = 2 }", "СЕМ-ДОЛЖЕН-БЫТЬ-ТИП"},

	{"модуль м; вход { пусть я: М.Т = 1 }", "СЕМ-НЕ-НАЙДЕН-МОДУЛЬ"},
	{"модуль м; конст М: Цел64 = 1; вход { пусть я: М.Т = 1 }", "СЕМ-ДОЛЖЕН-БЫТЬ-МОДУЛЬ"},
}

//===

func TestValid(t *testing.T) {
	fmt.Printf("--- valid tests: %d ---\n", len(valid_texts))
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

func TestErrors(t *testing.T) {
	fmt.Printf("--- tests for errors: %d ---\n", len(error_tests))
	t.Run("error tests", func(t *testing.T) {
		for _, e := range error_tests {
			checkForError(t, e.text, e.id)
		}
	})
}

func checkForError(t *testing.T, text, id string) {
	compile(text)
	if env.ErrorCount() == 0 {
		t.Errorf("An error is expected in text:\n%s\n", text)
		return
	}
	if id != "" {
		if env.GetErrorId(0) != id {
			t.Errorf("Expected '%s' error, got '%s' in text:\n%s\n", id, env.GetErrorId(0), text)
		}
	}
	env.ClearErrors()
}

func compile(text string) {
	var src = env.AddImmSource(text)

	m := parser.Parse(src)
	if env.ErrorCount() > 0 {
		return
	}
	Process(m)
}

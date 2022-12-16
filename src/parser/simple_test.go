package parser

import (
	"testing"
	"trivil/env"
)

var valid_texts = []string{
	"модуль м;",
	"модуль м",
	"модуль м\n",
	"модуль x; @внешняя фн print_int(a: Цел); вход { print_int(5) }",
	"модуль м; импорт 'путь'",
	"модуль м; пусть а: Цел",
	"модуль м; пусть а*: Цел",
}

var invalid_texts = []string{
	"иначе м",
	"модуль м; импорт;",
	"модуль м; импорт 12",
	"модуль м; пусть 'a': Цел",
	"модуль м; пусть 12: Цел",
	"модуль м; пусть a, б: Цел",
	"модуль м; пусть a+ Цел",
	"модуль м; пусть a: 12",
	//	"модуль м",
}

//===

func TestValid(t *testing.T) {
	t.Run("valid tests", func(t *testing.T) {
		for _, text := range valid_texts {
			checkValid(t, text)
		}
	})
}

func TestInvalid(t *testing.T) {
	t.Run("invalid tests", func(t *testing.T) {
		for _, text := range invalid_texts {
			checkInvalid(t, text)
		}
	})
}

func checkValid(t *testing.T, text string) {
	parseSrc(text)
	if env.ErrorCount() > 0 {
		t.Errorf("Unexpected %d errors in text:\n%s\n", env.ErrorCount(), text)
		env.ClearErrors()
	}
}

func checkInvalid(t *testing.T, text string) {
	parseSrc(text)
	if env.ErrorCount() > 0 {
		env.ClearErrors()
	} else {
		t.Errorf("Error(s) expected in text:\n%s\n", text)
	}
}

func parseSrc(text string) {
	var src = env.AddImmSource(text)

	Parse(src)
}

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
	"модуль м; конст а: Цел = 1",
	"модуль м; конст а*: Цел = 1",
	"модуль м; конст ( а: Цел = 1; б )",
	"модуль м; конст *( а: Цел = 1; б )",
	"модуль м; конст ( а: Цел = 1; б; в: Цел = 2; г )",
	"модуль м; тип а = []Цел",
	"модуль м; тип а* = []Цел",
	"модуль м; тип а = класс {}",
	"модуль м; тип а = класс (Цел) {}",
	"модуль м; тип а = класс { а: Цел }",
	"модуль м; тип а = класс { а: Цел; b: Строка }",
	"модуль м; тип а* = класс { а: Цел; b*: Строка }",

	"модуль м; фн Ф() {}",
	"модуль м; фн Ф(): Цел {}",
	"модуль м; фн Ф(а: Т) {}",
	"модуль м; фн Ф(а: Т, б: Цел) {}",
	"модуль м; фн Ф(а: Т, б: Цел): Строка {}",
	"модуль м; тип К = класс {}; фн (к: К) метод() {}",

	"модуль м; вход { a.b }",
	"модуль м; вход { a‹Цел› }",
	"модуль м; вход { a[] }",
	"модуль м; вход { a[1] }",
	"модуль м; вход { a[1, 2] }",
	"модуль м; вход { a[1, 2,] }",
	"модуль м; вход { a[1:2, 2:3] }",

	"модуль м; вход { a{} }",
	"модуль м; вход { a{б: 1} }",
	"модуль м; вход { a{б: 1,} }",
	"модуль м; вход { a{б: 1, в: 'фф'} }",
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
	"модуль м; конст а = 1",
	"модуль м; конст а: Цел",
	"модуль м; конст { а: Цел = 1 )",
	//	"модуль м",

	"модуль м; вход { a.1] }",
	"модуль м; вход { a‹› }",
	"модуль м; вход { a[1, 2:3] }",
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
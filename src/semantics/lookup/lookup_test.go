package lookup

import (
	"fmt"
	"testing"
	"trivil/env"
	"trivil/parser"
)

var valid_texts = []string{
	"модуль м; пусть я: Цел64; вход { я = 1 }",
	"модуль м; вход { пусть я: Цел64; я = 1 }",
	"модуль м; конст ц: Цел64 = 1; вход { пусть я: Цел64; я = ц }",
	"модуль м; фн Ф() {}; вход { Ф() }",
	"модуль м; тип Я = []Цел64; пусть я: Я; вход { }",
	"модуль м; тип Я = []Цел64; пусть я: Я; вход { я[0] := 1 }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я; я[0] := 1 }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я; я = Я[] }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я; пусть й: Цел64; я = Я[й] }",
	"модуль м; тип Я = []Цел64; вход { пусть я: Я; пусть й: Цел64; я = Я[й, й] }",
	"модуль м; тип Я = []Цел64; тип М = []Я; вход { }",
	"модуль м; тип М = []Я; тип Я = []Цел64; вход { }",
	"модуль м; тип К = класс {}; вход { }",
	"модуль м; тип К = класс { я: Цел64 }; вход { }",
	"модуль м; тип К = класс { я: Цел64; к: К }; вход { }",
	"модуль м; тип К = класс { я: Цел64; к: К }; вход { К{} }",
	"модуль м; тип К = класс { я: Цел64; б: Цел64 }; вход { К{я: 1} }",
	"модуль м; тип К = класс { я: Цел64; б: Цел64 }; вход { К{я: 1, б: 2} }",
}

var invalid_texts = []string{
	"иначе м",
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
	Process(m)
}

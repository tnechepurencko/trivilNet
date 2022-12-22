package check

import (
	"fmt"
	"testing"
	"trivil/env"
	"trivil/parser"
	"trivil/semantics/lookup"
)

type one struct {
	text string
	id   string
}

var error_tests = []one{
	{"модуль x; тип А = класс (Цел) {}", "СЕМ-БАЗА-НЕ-КЛАСС"},
	{"модуль x; тип А = класс { ц: Цел; ц: Цел}", "СЕМ-ДУБЛЬ-В-КЛАССЕ"},
	{"модуль x; тип А = класс { ц: Цел}; фн (а: А) ц() {}", "СЕМ-ДУБЛЬ-В-КЛАССЕ"},
	{"модуль x; тип А = класс {}; тип Б = класс (А){}; фн (а: А) Ф() {}; фн (б: Б) Ф(): Цел {}", "СЕМ-РАЗНЫЕ-ТИПЫ-МЕТОДОВ"},
	{"модуль x; тип А = класс {}; тип Б = класс (А){}; фн (а: А) Ф() {}; фн (б: Б) Ф(х: Цел) {}", "СЕМ-РАЗНЫЕ-ТИПЫ-МЕТОДОВ"},
	{"модуль x; тип А = класс {}; тип Б = класс (А){}; фн (а: А) Ф(х: Лог) {}; фн (б: Б) Ф(х: Цел) {}", "СЕМ-РАЗНЫЕ-ТИПЫ-МЕТОДОВ"},

	{"модуль x; вход { если 1 {} }", "СЕМ-ТИП-ВЫРАЖЕНИЯ"},
}

//===

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
	lookup.Process(m)
	if env.ErrorCount() > 0 {
		return
	}
	Process(m)
}
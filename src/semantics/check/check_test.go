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
	{"модуль x; тип А = класс (Цел64) {}", "СЕМ-БАЗА-НЕ-КЛАСС"},
	{"модуль x; тип А = класс { ц: Цел64; ц: Цел64}", "СЕМ-ДУБЛЬ-В-КЛАССЕ"},
	{"модуль x; тип А = класс { ц: Цел64}; фн (а: А) ц() {}", "СЕМ-ДУБЛЬ-В-КЛАССЕ"},
	{"модуль x; тип А = класс {}; тип Б = класс (А){}; фн (а: А) Ф() {}; фн (б: Б) Ф(): Цел64 {}", "СЕМ-РАЗНЫЕ-ТИПЫ-МЕТОДОВ"},
	{"модуль x; тип А = класс {}; тип Б = класс (А){}; фн (а: А) Ф() {}; фн (б: Б) Ф(х: Цел64) {}", "СЕМ-РАЗНЫЕ-ТИПЫ-МЕТОДОВ"},
	{"модуль x; тип А = класс {}; тип Б = класс (А){}; фн (а: А) Ф(х: Лог) {}; фн (б: Б) Ф(х: Цел64) {}", "СЕМ-РАЗНЫЕ-ТИПЫ-МЕТОДОВ"},

	{"модуль x; вход { если 1 {} }", "СЕМ-ТИП-ВЫРАЖЕНИЯ"},

	{"модуль x; вход { 1() }", "СЕМ-ВЫЗОВ-НЕ_ФУНКТИП"},

	{"модуль x; вход { пусть ц: Цел64 = ложь }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},
	{"модуль x; вход { пусть ц: Цел64 = 1; ц := ложь }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},
	{"модуль x; вход { пусть ц: Символ = '' }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},
	{"модуль x; вход { пусть ц: Символ = 'фф' }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},

	{"модуль x; фн Ф() {}; вход { Ф(1) }", "СЕМ-ЧИСЛО-АРГУМЕНТОВ"},
	{"модуль x; фн Ф(ц: Цел64) {}; вход { Ф() }", "СЕМ-ЧИСЛО-АРГУМЕНТОВ"},
	{"модуль x; фн Ф(ц1: Цел64, ц2: Цел64) {}; вход { Ф(1) }", "СЕМ-ЧИСЛО-АРГУМЕНТОВ"},
	{"модуль x; фн Ф(ц: Цел64) {}; вход { Ф(ложь) }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},
	{"модуль x; фн Ф(л: Лог) {}; вход { Ф(1) }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},

	{"модуль x; фн Ф(ц1: Цел64, ц2: ...Цел64) {}; вход { Ф() }", "СЕМ-ВАРИАДИК-ЧИСЛО-АРГУМЕНТОВ"},
	{"модуль x; фн Ф(ц1: Цел64, ц2: ...Цел64) {}; вход { Ф(ложь) }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},
	{"модуль x; фн Ф(ц1: Цел64, ц2: ...Цел64) {}; вход { Ф(1, ложь) }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},

	{"модуль x; вход { пусть ц = 0; ц = Цел64 }", "СЕМ-ТИП-В-ВЫРАЖЕНИИ"},
	{"модуль x; вход { пусть ц = 0; ц = 1 + Цел64 }", "СЕМ-ТИП-В-ВЫРАЖЕНИИ"},
	{"модуль x; вход { Лог.а }", "СЕМ-ТИП-В-ВЫРАЖЕНИИ"},

	{"модуль x; вход { если ~ 1 {} }", "СЕМ-ОШ-УНАРНАЯ-ТИП"},
	{"модуль x; вход { пусть ц = ложь; ц++ }", "СЕМ-ОШ-УНАРНАЯ-ТИП"},
	{"модуль x; вход { пусть ц = ложь; ц-- }", "СЕМ-ОШ-УНАРНАЯ-ТИП"},

	{"модуль x; вход { ложь + ложь }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},
	{"модуль x; вход { ложь - ложь }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},
	{"модуль x; вход { ложь * ложь }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},
	{"модуль x; вход { ложь / ложь }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},
	{"модуль x; вход { ложь % ложь }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},

	{"модуль x; вход { 1 & ложь }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},
	{"модуль x; вход { 1 | ложь }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},
	{"модуль x; вход { ложь & 1 }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},
	{"модуль x; вход { ложь | 1 }", "СЕМ-ОШ-ТИП-ОПЕРАНДА"},

	{"модуль x; вход { 1 + ложь }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 - ложь }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 * ложь }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 / ложь }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 % ложь }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},

	{"модуль x; вход { 1 + 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 - 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 * 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 / 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 % 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},

	{"модуль x; вход { 1 = 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 # 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 < 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 <= 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 > 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},
	{"модуль x; вход { 1 >= 1.0 }", "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ"},

	{"модуль x; вход { 1(:Лог) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; вход { ложь(:Лог) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; вход { ложь(:Цел64) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; вход { ложь(:Вещ64) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; вход { 1.0(:Символ) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; вход { 1(:Строка) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; вход { 1.0(:Строка) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; тип А = []Лог; вход { пусть а = А[]; пусть б = а(:Строка) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},
	{"модуль x; тип А = []Лог; вход { пусть а = 'при'(:А) }", "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА"},

	{"модуль x; вход { 1(:Цел64) }", "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ"},
	{"модуль x; вход { 1.0(:Вещ64) }", "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ"},
	{"модуль x; вход { 1(:Символ)(:Символ) }", "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ"},
	{"модуль x; вход { 'при'(:Строка) }", "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ"},

	{"модуль x; вход { 256(:Байт) }", "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ"},
	{"модуль x; вход { 'ф'(:Байт) }", "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ"},
	{"модуль x; вход { 'фи'(:Байт) }", "СЕМ-ДЛИНА-СТРОКИ-НЕ-1"},
	{"модуль x; вход { ''(:Байт) }", "СЕМ-ДЛИНА-СТРОКИ-НЕ-1"},

	{"модуль x; тип А = класс {}; тип Б = класс {}; вход { А{}(:Б) }", "СЕМ-ДОЛЖЕН-БЫТЬ-НАСЛЕДНИКОМ"},
	{"модуль x; тип А = класс (Б) {}; тип Б = класс {}; вход { А{}(:Б) }", "СЕМ-ДОЛЖЕН-БЫТЬ-НАСЛЕДНИКОМ"},

	{"модуль x; фн Ф() { вернуть 1 }", "СЕМ-ОШ-ВЕРНУТЬ-ЛИШНЕЕ"},
	{"модуль x; фн Ф(): Цел64 { вернуть }", "СЕМ-ОШ-ВЕРНУТЬ-НУЖНО"},
	{"модуль x; фн Ф(): Цел64 { вернуть ложь }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},
	{"модуль x; вход { прервать }", "СЕМ-ПРЕРВАТЬ-ВНЕ-ЦИКЛА"},

	{"модуль x; вход { надо ложь иначе {} }", "СЕМ-НЕ-ЗАВЕРШАЮЩИЙ"},
	{"модуль x; вход { надо ложь иначе { пусть ц = 0 } }", "СЕМ-НЕ-ЗАВЕРШАЮЩИЙ"},

	{"модуль x; вход { пусть а = ложь[] }", "СЕМ-КОМПОЗИТ-НЕТ-ТИПА"},
	{"модуль x; вход { пусть а = Лог[] }", "СЕМ-МАССИВ-КОМПОЗИТ-ОШ-ТИП"},
	{"модуль x; тип А = []Цел64; вход { пусть а = А[ложь: 1] }", "СЕМ-МАССИВ-КОМПОЗИТ-ТИП-КЛЮЧА"},
	{"модуль x; тип А = []Цел64; вход { пусть ц = 0; пусть а = А[ц: 1] }", "СЕМ-ОШ-КОНСТ-ВЫРАЖЕНИЕ"},
	{"модуль x; тип А = []Цел64; вход { пусть а = А[1: ложь] }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},

	{"модуль x; вход { пусть а = ложь{} }", "СЕМ-КОМПОЗИТ-НЕТ-ТИПА"},
	{"модуль x; вход { пусть а = Лог{} }", "СЕМ-КЛАСС-КОМПОЗИТ-ОШ-ТИП"},
	{"модуль x; тип К = класс { ц: Цел64}; вход { пусть к = К{я: 1}}", "СЕМ-КЛАСС-КОМПОЗИТ-НЕТ-ПОЛЯ"},
	{"модуль x; тип К = класс { ц: Цел64}; фн (к: К) я() {}; вход { пусть к = К{я: 1}}", "СЕМ-КЛАСС-КОМПОЗИТ-НЕ-ПОЛE"},
	{"модуль x; тип К = класс { ц: Цел64}; вход { пусть к = К{ц: ложь}}", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},
	{"модуль x; вход { пусть ц = 0; ц.п }", "СЕМ-ОЖИДАЛСЯ-ТИП-КЛАССА"},
	{"модуль x; тип К = класс {}; вход { пусть к = К{}; к.ц }", "СЕМ-ОЖИДАЛОСЬ-ПОЛЕ-ИЛИ-МЕТОД"},
	{"модуль x; тип К = класс { ц: Цел64}; вход { пусть к = К{}; к.ц := ложь }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},

	{"модуль x; вход { пусть ц = 0; ц[1] }", "СЕМ-ОЖИДАЛСЯ-ТИП-МАССИВА"},
	{"модуль x; тип А = []Цел64; вход { пусть а = А[]; а[ложь] }", "СЕМ-ОШ-ТИП-ИНДЕКСА"},

	{"модуль x; вход { 2 := 1 }", "СЕМ-НЕ-ПРИСВОИТЬ"},
	{"модуль x; вход { 2++ }", "СЕМ-НЕ-ПРИСВОИТЬ"},
	{"модуль x; вход { 2-- }", "СЕМ-НЕ-ПРИСВОИТЬ"},
	{"модуль x; конст к: Цел64 = 1; вход { к := 1 }", "СЕМ-НЕ-ПРИСВОИТЬ"},
	{"модуль x; конст к: Цел64 = 1; вход { к++ }", "СЕМ-НЕ-ПРИСВОИТЬ"},
	{"модуль x; конст к: Цел64 = 1; вход { к := 1 }", "СЕМ-НЕ-ПРИСВОИТЬ"},
	{"модуль x; вход { пусть ц = 1; ц := 2 }", "СЕМ-НЕ-ПРИСВОИТЬ"},

	{"модуль x; вход { 1 }", "СЕМ-ЗНАЧЕНИЕ-НЕ-ИСПОЛЬЗУЕТСЯ"},
	{"модуль x; вход { пусть ц = 1; ц }", "СЕМ-ЗНАЧЕНИЕ-НЕ-ИСПОЛЬЗУЕТСЯ"},
	{"модуль x; вход { пусть ц = 1; ц = 2}", "СЕМ-ЗНАЧЕНИЕ-НЕ-ИСПОЛЬЗУЕТСЯ"},

	{"модуль x; вход { пусть а = длина() }", "СЕМ-ОШ-ЧИСЛО-АРГ-СТДФУНК"},
	{"модуль x; вход { пусть а = длина(1, 2) }", "СЕМ-ОШ-ЧИСЛО-АРГ-СТДФУНК"},
	{"модуль x; вход { пусть а = длина(1) }", "СЕМ-ДЛИНА-ОШ-ТИП-АРГ"},
	{"модуль x; тип А = []Цел64; вход { пусть а = А[]; пусть л: Лог = длина(а) }", "СЕМ-НЕСОВМЕСТИМО-ПРИСВ"},

	{"модуль x; вход { авария(1) }", "СЕМ-ТИП-ВЫРАЖЕНИЯ"},
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

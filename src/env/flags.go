package env

import (
	"flag"
)

var (
	JustLexer = flag.Bool("lexer", false, "только лексический анализ - показать токены")
	DoTrace   = flag.Bool("trace", false, "включить трассировку парсера")
	DoGen     = flag.Bool("gen", true, "включить генерацию")
	ShowAST   = flag.Int("ast", 0, "ast=1 - after parser; ast=2 - after analyzer")
)

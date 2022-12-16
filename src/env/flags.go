package env

import (
	"flag"
)

var (
	JustLexer = flag.Bool("lexer", false, "только лексический анализ - показать токены")
	TraceFlag = flag.Bool("trace", false, "включить трассировку парсера")
)

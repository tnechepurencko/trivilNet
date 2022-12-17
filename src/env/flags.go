package env

import (
	"flag"
)

var (
	JustLexer = flag.Bool("lexer", false, "только лексический анализ - показать токены")
	TraceFlag = flag.Bool("trace", false, "включить трассировку парсера")
	ShowAST   = flag.Int("ast", 0, "ast=1 - after parser; ast=2 - after analyzer")
)

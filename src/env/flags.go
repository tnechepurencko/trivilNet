package env

import (
	"flag"
)

var (
	JustLexer    = flag.Bool("lexer", false, "только лексический анализ - показать токены")
	TraceParser  = flag.Bool("trace_parser", false, "включить трассировку парсера")
	TraceCompile = flag.Bool("trace_compile", false, "включить трассировку компиляции программы")
	DoGen        = flag.Bool("gen", true, "включить генерацию")
	ShowAST      = flag.Int("ast", 0, "ast=1 - after parser; ast=2 - after analyzer")
)

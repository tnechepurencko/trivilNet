package env

import (
	"flag"
)

var (
	TraceFlag = flag.Bool("trace", false, "включить трассировку парсера")
)

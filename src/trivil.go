package main

import (
	"flag"
	"fmt"
	"os"

	"trivil/compiler"
	"trivil/env"
	"trivil/lexer"
)

func main() {
	// флаги определены в env.flags
	flag.Parse()
	spath := flag.Arg(0)
	if spath == "" {
		fmt.Println("Использование: tric (folder | file.tri)")
		os.Exit(1)
	}

	fmt.Println("Тривиль-0 компилятор v0.0")
	env.Init()

	//fmt.Printf("%v\n", src.Bytes)
	if *env.JustLexer {
		testLexer(spath)
	} else {
		compiler.Compile(spath)
	}

	env.ShowErrors()

	if env.ErrorCount() == 0 {
		fmt.Printf("Без ошибок\n")
	} else {
		os.Exit(1)
	}

}

func testLexer(arg string) {

	var files = env.GetSources(arg)
	var src = files[0]
	if src.Err != nil {
		fmt.Printf("Ошибка чтения исходного файла '%s': %s\n", arg, src.Err.Error())
		os.Exit(1)
	}

	var lex = new(lexer.Lexer)
	lex.Init(src)

	for true {
		pos, tok, lit := lex.Scan()
		if tok == lexer.EOF {
			break
		}

		_, line, col := env.SourcePos(pos)
		if lit == "" {
			fmt.Printf("%d:%d %v\n", line, col, tok)
		} else {
			fmt.Printf("%d:%d %v '%s'\n", line, col, tok, lit)
		}
	}
}

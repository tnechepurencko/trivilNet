package main

import (
	"flag"
	"fmt"
	"os"

	"trivil/ast"
	"trivil/env"
	"trivil/genc"
	"trivil/lexer"
	"trivil/parser"
	"trivil/semantics"
)

func main() {
	// флаги определены в env.flags
	flag.Parse()
	arg := flag.Arg(0)
	if arg == "" {
		fmt.Println("Использование: trivil name.tri")
		os.Exit(1)
	}

	fmt.Println("Тривиль-0 компилятор v0.0")
	env.Init()

	src := env.AddSource(arg)
	if src.Err != nil {
		fmt.Printf("Ошибка чтения исходного файла '%s': %s\n", arg, src.Err.Error())
		os.Exit(1)
	}

	//fmt.Printf("%v\n", src.Bytes)
	if *env.JustLexer {
		testLexer(src)
	} else {
		compile(src)
	}

	env.ShowErrors()

	if env.ErrorCount() == 0 {
		fmt.Printf("Без ошибок\n")
	}

}

func compile(src *env.Source) {
	var m = parser.Parse(src)
	if env.ErrorCount() != 0 {
		return
	}

	if *env.ShowAST >= 1 {
		fmt.Println(ast.SExpr(m))
	}

	semantics.Analyse(m)
	if env.ErrorCount() != 0 {
		return
	}

	if *env.ShowAST >= 2 {
		fmt.Println(ast.SExpr(m))
	}

	if *env.DoGen {
		genc.Generate(m)
	}

}

func testLexer(src *env.Source) {
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

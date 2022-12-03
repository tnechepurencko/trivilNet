package main

import (
	"flag"
	"fmt"
	"os"
	"trivil/env"
	"trivil/lexer"
)

func main() {
	flag.Parse()
	arg := flag.Arg(0)
	if arg == "" {
		fmt.Println("Использование: trivil name.tri")
		os.Exit(1)
	}

	fmt.Println("Тривиль-0 компилятор v0.0")
	env.Init()

	src := env.ReadSource(arg)
	if src.Err != nil {
		fmt.Printf("Ошибка чтения исходного файла '%s': %s\n", arg, src.Err.Error())
		os.Exit(1)
	}

	//fmt.Printf("%v\n", src.Bytes)

	var lex = new(lexer.Lexer)
	lex.Init(src)

	for true {
		pos, tok, lit := lex.Scan()
		if tok == lexer.EOF {
			break
		}
		fmt.Printf("%d %v %s\n", pos, tok, lit)
	}

	env.ShowErrors()

}

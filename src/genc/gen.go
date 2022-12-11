package genc

import (
	"fmt"
	"trivil/ast"
)

var _ = fmt.Printf

type genContext struct {
	module *ast.Module
	hlines []string
	clines []string
}

func Generate(m *ast.Module) {
	fmt.Printf("genc\n")

	var genc = &genContext{
		module: m,
		hlines: make([]string, 0),
		clines: make([]string, 0),
	}

	genc.startCode()
	genc.genModule()
	genc.finishCode()

	genc.show()
}

func (genc *genContext) h(format string, args ...interface{}) {
	genc.hlines = append(genc.hlines, fmt.Sprintf(format, args...))
}

func (genc *genContext) c(format string, args ...interface{}) {
	genc.clines = append(genc.clines, fmt.Sprintf(format, args...))
}

func (genc *genContext) startCode() {
	var hname = fmt.Sprintf("_%s_H", genc.module.Name)
	genc.h("#ifndef %s", hname)
	genc.h("#define %s", hname)

	genc.c("#include \"%s\"", "header")

}

func (genc *genContext) finishCode() {
	genc.h("#endif")
}

//====

func (genc *genContext) show() {
	fmt.Println("---- header ----")
	for _, s := range genc.hlines {
		fmt.Println(s)
	}
	fmt.Println("---- c code ----")
	for _, s := range genc.clines {
		fmt.Println(s)
	}
	fmt.Println("---- end c  ----")
}

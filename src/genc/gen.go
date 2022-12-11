package genc

import (
	"fmt"
	"os"
	"path"
	"strings"

	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

type genContext struct {
	module  *ast.Module
	outname string
	hlines  []string
	clines  []string
}

func Generate(m *ast.Module) {

	var genc = &genContext{
		module:  m,
		outname: env.OutName(m.Name),
		hlines:  make([]string, 0),
		clines:  make([]string, 0),
	}

	genc.startCode()
	genc.genModule()
	genc.finishCode()

	//genc.show()
	genc.save()

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

	genc.c("#include \"trirun.h\"")
	genc.c("#include \"%s\"", genc.outname+".h")

}

func (genc *genContext) finishCode() {
	genc.h("#endif")
}

//====

func (genc *genContext) save() {
	var folder = env.PrepareOutFolder()

	writeFile(folder, genc.outname, ".h", genc.hlines)
	writeFile(folder, genc.outname, ".c", genc.clines)
}

func writeFile(folder, name, ext string, lines []string) {

	var filename = path.Join(folder, name+ext)

	var out = strings.Join(lines, "\n")

	var err = os.WriteFile(filename, []byte(out), 0755)

	if err != nil {
		panic("Ошибка записи файла " + filename + ": " + err.Error())
	}
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

package genc

import (
	"fmt"
	"os"

	//	"path"
	//	"strings"
	"runtime"

	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func BuildExe(modules []*ast.Module) {
	fmt.Printf("build: %s\n", runtime.GOOS)

	var template = findTemplate("exe-" + runtime.GOOS)
	if template == "" {
		return
	}
}

func findTemplate(name string) string {

	_, err := os.ReadFile("conf_genc.txt")
	if err != nil {
		env.AddProgramError("ГЕН-ОШ-КОНФ-ФАЙЛА", err.Error())
		return ""
	}

	//	var lines = strings.Split(string(buf[:]), "\n")

	return ""
}

//====

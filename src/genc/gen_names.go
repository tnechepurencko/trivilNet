package genc

import (
	"fmt"
	//	"trivil/ast"
)

var _ = fmt.Printf

// prefixes for generated names
const (
	nm_stringLiteral = "strlit"
)

// run-time API
const (
	rt_prefix           = "tri_"
	rt_newLiteralString = rt_prefix + "newLiteralString"
	rt_lenString        = rt_prefix + "lenString"
)

func (genc *genContext) localName(prefix string) string {
	genc.autoNo++
	return fmt.Sprintf("%s%d", prefix, genc.autoNo)
}

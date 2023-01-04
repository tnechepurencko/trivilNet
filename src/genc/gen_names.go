package genc

import (
	"fmt"
	//	"trivil/ast"
)

var _ = fmt.Printf

const typeNamePrefix = "T"

// meta информация
const (
	nm_VT_suffix       = "_VT"
	nm_VT_field        = "_vtable_"
	nm_meta_suffix     = "_Meta"
	nm_meta_field      = "_meta_"
	nm_meta_var_prefix = "meta_"
)

// prefixes for generated names
const (
	nm_stringLiteral = "strlit"
)

// run-time API
const (
	rt_prefix = "tri_"

	rt_newLiteralString = rt_prefix + "newLiteralString"
	rt_lenString        = rt_prefix + "lenString"

	rt_newVector = rt_prefix + "newVector"
	rt_lenVector = rt_prefix + "lenVector"
	rt_vcheck    = rt_prefix + "vcheck"

	rt_newObject = rt_prefix + "newObject"
)

func (genc *genContext) localName(prefix string) string {
	genc.autoNo++
	return fmt.Sprintf("%s%d", prefix, genc.autoNo)
}

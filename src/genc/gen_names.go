package genc

import (
	"fmt"
	//	"trivil/ast"
)

var _ = fmt.Printf

const typeNamePrefix = "T"

// класс струкура и мета информация
const (
	nm_class_struct_suffix = "_ST"
	nm_class_fields        = "f"
	nm_base_fields         = "_B"
	nm_VT_field            = "vtable"

	nm_VT_suffix       = "_VT"
	nm_meta_suffix     = "_Meta"
	nm_meta_field      = "_meta_"
	nm_desc_var_suffix = "_desc"
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

	rt_newObject      = rt_prefix + "newObject"
	rt_checkClassType = rt_prefix + "checkClassType"

	rt_convert = rt_prefix
)

func (genc *genContext) localName(prefix string) string {
	genc.autoNo++
	return fmt.Sprintf("%s%d", prefix, genc.autoNo)
}

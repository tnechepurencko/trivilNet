package genc

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
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

	rt_crash = rt_prefix + "crash"
)

func (genc *genContext) localName(prefix string) string {
	genc.autoNo++
	return fmt.Sprintf("%s%d", prefix, genc.autoNo)
}

//====

func (genc *genContext) declName(d ast.Decl) string {

	out, ok := genc.declNames[d]
	if ok {
		return out
	}

	f, is_fn := d.(*ast.Function)

	if is_fn && f.External {
		var name = f.ExternalName
		if name == "" {
			name = env.OutName(f.Name)
		}
		genc.declNames[d] = name

		return name
	}

	out = ""
	var host = d.GetHost()
	if host != nil {
		out = genc.declName(host) + "__"
	}

	var prefix = ""
	if _, ok := d.(*ast.TypeDecl); ok {
		prefix = typeNamePrefix
	}

	out += prefix + env.OutName(d.GetName())

	genc.declNames[d] = out

	return out
}

func (genc *genContext) outName(name string) string {
	return env.OutName(name)
}

func (genc *genContext) functionName(f *ast.Function) string {

	if f.Recv != nil {
		return genc.typeRef(f.Recv.Typ) + "_" + genc.outName(f.Name)
	}
	return genc.declName(f)
}

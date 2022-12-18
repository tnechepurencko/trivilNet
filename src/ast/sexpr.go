package ast

import (
	"fmt"
	"reflect"
	"strings"
)

var _ = fmt.Printf

func SExpr(n interface{}) string {
	return sexpr(reflect.ValueOf(n))
}

func sexpr(v reflect.Value) string {

	v, ok := getStruct(v)
	if !ok {
		return "_" + v.Type().Kind().String()
	}

	if v.Type().Name() == "Scope" {
		return ""
	}

	var fs = ""
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		if v.Type().Field(i).Name == "Pos" {
			continue
		}

		//fmt.Println(f.Type().String())

		switch f.Kind() {
		case reflect.Int:
			fs += fmt.Sprintf(" (%s %d)", v.Type().Field(i).Name, f.Int())
		case reflect.Bool:
			if f.Bool() {
				fs += " " + v.Type().Field(i).Name
			}
		case reflect.String:
			fs += " \"" + f.String() + "\""
		case reflect.Pointer:
			if !f.IsNil() {
				fs += " " + sexpr(f)
			}
		case reflect.Slice:
			var list = slice(f)

			fs += fmt.Sprintf(" [%s]", strings.Join(list, " "))
		case reflect.Interface:
			if !f.IsNil() {
				fs += " " + sexpr(f.Elem())
			}
		case reflect.Struct:
			sname := f.Type().Name()
			if sname == "DeclBase" {
				name := f.FieldByName("Name")
				fs += " \"" + name.String() + "\""
				exported := f.FieldByName("Exported")
				if exported.Bool() {
					fs += " Exported"
				}
			} else if strings.HasSuffix(sname, "Base") {
				// игнорирую
			} else {
				fs += " " + sexpr(f)
			}
		}
	}
	var str = fmt.Sprintf("(%s%s)", v.Type().Name(), fs)

	return str
}

func slice(v reflect.Value) []string {
	var s = make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		s[i] = sexpr(v.Index(i))
	}

	return s
}

func getStruct(v reflect.Value) (reflect.Value, bool) {

	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		return v, true
	}
	return v, false
}

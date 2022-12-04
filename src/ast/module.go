package ast

import (
	"fmt"
	//	"trivil/env"
)

var _ = fmt.Printf

//=== модуль

type Module struct {
	//source *env.Source
	Name  string
	Inner *Scope
	Entry *EntryFn
}

//=== вход

type EntryFn struct {
	//source *env.Source
}

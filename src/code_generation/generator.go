package code_generation

import "trivil/compiler"
import "trivil/ast"

type Code struct {
}

func StartGeneration(cc *compiler.CompileContext) *Code {
	var c = &Code{}

	for _, imps := range cc.Main.Imports {
		c.GenerateModule(imps.Mod)
	}

	for _, decl := range cc.Main.Decls {
		c.GenerateDeclaration(decl)
	}

	for _, stmt := range cc.Main.Entry.Seq.Statements {
		c.GenerateStatement(stmt)
	}

	return c
}

func (c *Code) GenerateDeclaration(decl ast.Decl) {
	switch x := decl.(type) {
	case *ast.TypeDecl:
		c.GenerateTypeDecl(x)
	case *ast.Function:
		c.GenerateFunction(x)
	}
}

func (c *Code) GenerateStatement(stmt ast.Statement) {
	switch x := stmt.(type) {
	case *ast.ExprStatement:
		c.GenerateExpr(x.X)
	}
}

func (c *Code) GenerateModule(td *ast.Module) {

}

func (c *Code) GenerateTypeDecl(td *ast.TypeDecl) {

}

func (c *Code) GenerateFunction(td *ast.Function) {

}

func (c *Code) GenerateExpr(expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.CallExpr:
		c.GenerateCallExpr(x)
	}
}

func (c *Code) GenerateCallExpr(td *ast.CallExpr) {

}

package code_generation

import "trivil/compiler"
import "trivil/ast"

type Code struct {
}

func StartGeneration(cc *compiler.CompileContext) *Code {
	var c = &Code{}

	c.GenerateModule(cc.Main)

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
	case *ast.Guard:
		c.GenerateIf(x)
	}
}

func (c *Code) GenerateIf(ifStmt *ast.Guard) {

}

func (c *Code) GenerateModule(module *ast.Module) {
	for _, imps := range module.Imports {
		c.GenerateModule(imps.Mod)
	}

	for _, decl := range module.Decls {
		c.GenerateDeclaration(decl)
	}

	for _, stmt := range module.Entry.Seq.Statements {
		c.GenerateStatement(stmt)
	}
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

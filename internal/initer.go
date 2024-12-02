package internal

import (
	"go/ast"
	"go/token"
)

type InitFuncBuilder struct {
	varToFunc map[string]string
}

func NewInitBuilder() *InitFuncBuilder {
	return &InitFuncBuilder{
		varToFunc: make(map[string]string),
	}
}

func (i *InitFuncBuilder) AddInit(varName, funcName string) {
	i.varToFunc[varName] = funcName
}

func (i *InitFuncBuilder) GenInit() *ast.FuncDecl {
	var initFuncDecl = &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}

	for varName, funcName := range i.varToFunc {
		var assignStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(varName)},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent(funcName)},
		}
		initFuncDecl.Body.List = append(initFuncDecl.Body.List, assignStmt)
	}
	return initFuncDecl
}

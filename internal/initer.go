package internal

import (
	"go/ast"
	"go/token"
)

type InitBuilder struct {
	varToFunc map[string]string
}

func NewInitBuilder() *InitBuilder {
	return &InitBuilder{
		varToFunc: make(map[string]string),
	}
}

func (i *InitBuilder) AddInit(varName, funcName string) {
	i.varToFunc[varName] = funcName
}

func (i *InitBuilder) GenInit() *ast.FuncDecl {
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

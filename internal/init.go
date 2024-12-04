package internal

import (
	"go/ast"
	"go/token"
)

type InitFuncBuilder struct {
	varToFunc map[string]string
}

func NewInitFuncBuilder() *InitFuncBuilder {
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
	var names = make([]string, 0, len(i.varToFunc))
	for k, _ := range i.varToFunc {
		names = append(names, k)
	}
	for _, varName := range names {
		var assignStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(varName)},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent(i.varToFunc[varName])},
		}
		initFuncDecl.Body.List = append(initFuncDecl.Body.List, assignStmt)
	}
	return initFuncDecl
}

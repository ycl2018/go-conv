package internal

import (
	"go/ast"
	"go/token"
)

func buildIfStmt(lh string, op token.Token, rh string) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent(lh),
			Op: op,
			Y:  ast.NewIdent(rh),
		},
		Body: &ast.BlockStmt{},
	}
}

func buildAssignStmt(lh string, rh string) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(lh)},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ast.NewIdent(rh)},
	}
}

func buildVarDecl(varName string, typeName string) *ast.DeclStmt {
	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{ast.NewIdent(varName)},
					Type:  ast.NewIdent(typeName),
				},
			},
		}}
}

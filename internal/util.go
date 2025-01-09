package internal

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"unicode"
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

func buildDefineStmt(lh string, rh string) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(lh)},
		Tok: token.DEFINE,
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

var (
	starPattern    = regexp.MustCompile(`(Ptr)+(\w)`)
	bracketPattern = regexp.MustCompile(`\[(\d+)]`)
	slicePattern   = regexp.MustCompile(`(Slice)+(\w)`)
)

func cleanName(name string) string {
	name = strings.ReplaceAll(name, "*", "Ptr")
	name = strings.ReplaceAll(name, "[]", "Slice")
	// replace [number] with Array[number]
	name = bracketPattern.ReplaceAllStringFunc(name, func(s string) string {
		number := s[1 : len(s)-1]
		return "Array" + number
	})

	// replace non-letter characters
	var result strings.Builder
	var prevIsSpecial bool
	var first = true
	for _, char := range name {
		if unicode.IsLetter(char) || (unicode.IsDigit(char) && !first) {
			if prevIsSpecial || first {
				result.WriteRune(unicode.ToUpper(char))
				prevIsSpecial = false
			} else {
				result.WriteRune(char)
			}
			first = false
		} else if !unicode.IsSpace(char) {
			prevIsSpecial = true
		}
	}
	output := result.String()
	// replace ptr* to upper
	output = starPattern.ReplaceAllStringFunc(output, func(s string) string {
		return s[:len(s)-1] + strings.ToUpper(s[len(s)-1:])
	})
	// replace slice* to upper
	output = slicePattern.ReplaceAllStringFunc(output, func(s string) string {
		return s[:len(s)-1] + strings.ToUpper(s[len(s)-1:])
	})

	return output
}

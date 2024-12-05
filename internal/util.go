package internal

import (
	"fmt"
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

func buildCommentExpr(comment string) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "// " + comment,
		},
	}
}

var (
	starPattern    = regexp.MustCompile(`\*(\w)`)
	bracketPattern = regexp.MustCompile(`\[(\d+)](\w)`)
	slicePattern   = regexp.MustCompile(`\[](\w)`)
)

func cleanName(name string) string {
	// 替换 * 为 Ptr 并将其后的第一个字母转为大写
	name = starPattern.ReplaceAllStringFunc(name, func(s string) string {
		return "Ptr" + strings.ToUpper(s[1:])
	})

	// 替换 [number] 为 Array 并将其后的第一个字母转为大写
	name = bracketPattern.ReplaceAllStringFunc(name, func(s string) string {
		number := s[1 : len(s)-2] // 提取数字部分
		nextChar := s[len(s)-1:]  // 提取下一个字符
		return fmt.Sprintf("Array%s%s", number, strings.ToUpper(nextChar))
	})

	// 替换 [] 为 Slice 并将其后的第一个字母转为大写
	name = slicePattern.ReplaceAllStringFunc(name, func(s string) string {
		return "Slice" + strings.ToUpper(s[len(s)-1:])
	})

	// 去除其他非字符和数字的符号，并处理特殊符号后面的字母转为大写
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

	// 构建结果字符串
	output := result.String()
	return output
}

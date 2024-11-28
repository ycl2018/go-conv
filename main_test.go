package go_conv

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

type MethodAndDoc struct {
	Method  *types.Func
	AstFunc *ast.Field
	Doc     *ast.CommentGroup
}

type Var struct {
	ImportType string
	Name       string
	Type       *types.Var
}

type Func struct {
	Doc        []string
	Name       string
	Src        *Var
	Dest       *Var
	Assignment []string
}

func Test(t *testing.T) {
	inputFile := "testdata/setup.go"
	absPath, err := filepath.Abs(inputFile)
	if err != nil {
		t.Fatalf("abs err:%v", err)
	}
	_, err = os.Stat(inputFile)
	if err != nil {
		t.Fatalf("inputFile:%s err:%v", inputFile, err)
	}
	const parserLoadMode = packages.NeedName | packages.NeedImports | packages.NeedDeps |
		packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo

	var srcAstFile *ast.File
	loadConf := &packages.Config{
		Mode: parserLoadMode,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			//log.Printf("parse filename:%s", filename)
			if filename == absPath {
				r, e := parser.ParseFile(fset, filename, src, parser.SkipObjectResolution|parser.ParseComments)
				srcAstFile = r
				return r, e
			} else {
				return parser.ParseFile(fset, filename, src, parser.SkipObjectResolution)
			}
		},
	}

	initial, err := packages.Load(loadConf, "file="+inputFile)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(initial)
	//t.Log(srcAstFile)
	srcPackage := initial[0]
	var methods []MethodAndDoc

	var importName = map[string]string{} // path to name
	for _, spec := range srcAstFile.Imports {
		var name string
		if spec.Name == nil {
			index := strings.LastIndex(spec.Path.Value, "/")
			name = spec.Path.Value[index+1:]
		} else {
			name = spec.Name.Name
			if name == "." {
				name = ""
			}
		}
		importName[spec.Path.Value] = name
	}
	// 当前包：
	importName[srcPackage.PkgPath] = ""

	for _, decl := range srcAstFile.Decls {
		// we only resolve the package level Decl
		if gd, ok := decl.(*ast.GenDecl); ok {
			var isTarget = false
			if gd.Doc == nil {
				continue
			}
			for _, comment := range gd.Doc.List {
				if strings.Contains(comment.Text, ":convergen") {
					isTarget = true
					break
				}
			}
			if !isTarget {
				continue
			}
			for _, spec := range gd.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					tye := srcPackage.TypesInfo.TypeOf(ts.Type)
					if itye, ok := tye.(*types.Interface); ok {
						astMethods := ts.Type.(*ast.InterfaceType).Methods.List
						for i := range itye.NumMethods() {
							method := itye.Method(i)
							var md = MethodAndDoc{
								Method:  method,
								Doc:     astMethods[i].Doc,
								AstFunc: astMethods[i],
							}
							methods = append(methods, md)
						}
					}
				}
			}
		}
	}
	var f = &ast.File{
		Doc:     nil,
		Package: srcAstFile.Package,
		Name: &ast.Ident{
			NamePos: 0,
			Name:    "test_gen.go",
			Obj:     nil,
		},
		Imports: srcAstFile.Imports,
	}
	// parse method
	for _, method := range methods {
		buildFunc(f, srcPackage.Types, importName, &method)
	}
	var sb strings.Builder
	err = printer.Fprint(&sb, token.NewFileSet(), f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(sb.String())
}

func buildFunc(f *ast.File, pkg *types.Package, importType map[string]string, method *MethodAndDoc) {
	src := method.Method.Signature().Params().At(0)
	dest := method.Method.Signature().Results().At(0)
	funcName := method.AstFunc.Names[0].Name
	// add a func
	fn := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: funcName,
		},
		Type: method.AstFunc.Type.(*ast.FuncType),
		Body: &ast.BlockStmt{},
	}
	// todo: check fn Type 参数是否有命名
	for _, field := range fn.Type.Params.List {
		if len(field.Names) == 0 {
			field.Names = []*ast.Ident{ast.NewIdent("src")}
			src = types.NewVar(0, src.Pkg(), "src", src.Type())
		}
	}
	for _, field := range fn.Type.Results.List {
		if len(field.Names) == 0 {
			field.Names = []*ast.Ident{ast.NewIdent("dest")}
			dest = types.NewVar(0, dest.Pkg(), "dest", dest.Type())
		}
	}
	f.Decls = append(f.Decls, fn)
	builder := &Builder{Package: pkg, ImportName: importType}
	stmts := builder.buildStmt(dest, src)
	fn.Body.List = append(fn.Body.List, stmts...)
	fn.Body.List = append(fn.Body.List, &ast.ExprStmt{X: ast.NewIdent("return")})
}

type Builder struct {
	Package    *types.Package
	ImportName map[string]string
}

func (b *Builder) buildStmt(dst *types.Var, src *types.Var) []ast.Stmt {
	if dst == nil {
		return nil
	}
	var stmts []ast.Stmt
	switch dstType := dst.Type().(type) {
	case *types.Pointer:
		if _, ok := src.Type().(*types.Pointer); !ok {
			log.Fatalf("src type is not a pointer:%s", src.String())
			return nil
		}
		ifStmt := &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: src.Name(),
				},
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			},
			Body: &ast.BlockStmt{},
		}
		// dst = new(dst.Type)
		destPtr, srcPtr := dst.Type().(*types.Pointer), src.Type().(*types.Pointer)
		dstElemVar := types.NewVar(0, b.Package, "*"+dst.Name(), destPtr.Elem())
		srcElemVar := types.NewVar(0, b.Package, "*"+src.Name(), srcPtr.Elem())
		initAssign := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{
				Fun:  ast.NewIdent("new"),
				Args: []ast.Expr{ast.NewIdent(importType(b.ImportName, dstElemVar))},
			}},
		}
		ifStmt.Body.List = append(ifStmt.Body.List, initAssign)
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		ifStmt.Body.List = append(ifStmt.Body.List, elementStmt...)
		stmts = append(stmts, ifStmt)
		return stmts
	case *types.Struct:
		srcType, ok := src.Type().(*types.Struct)
		if !ok {
			log.Fatalf("src type is not a struct:%s", src.String())
			return nil
		}
		srcName := strings.TrimPrefix(src.Name(), "*")
		dstName := strings.TrimPrefix(dst.Name(), "*")
		for i := range dstType.NumFields() {
			dstField := dstType.Field(i)
			dstFieldName := dstField.Name()
			// match srcField
			var match bool
			for j := range srcType.NumFields() {
				srcField := srcType.Field(j)
				if srcField.Name() == dstFieldName {
					dstVarName := dstName + "." + dstFieldName
					srcVarName := srcName + "." + srcField.Name()
					dstVar := types.NewVar(0, b.Package, dstVarName, dstField.Type())
					srcVar := types.NewVar(0, b.Package, srcVarName, srcField.Type())
					fieldStmt := b.buildStmt(dstVar, srcVar)
					stmts = append(stmts, fieldStmt...)
					match = true
					break
				}
			}
			if !match {
				log.Fatalf("src field %s not found in struct:%s", dstFieldName, srcType.String())
			}
		}
		return stmts
	case *types.Array:
		srcType, ok := src.Type().(*types.Array)
		if !ok {
			log.Fatalf("src type is not a array:%s", src.String())
			return nil
		}
		if srcType.Len() != dstType.Len() {
			log.Fatalf("src array len is not equal with dst")
			return nil
		}
		// for i := 0; i<n; i++ {}
		forStmt := &ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("i")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{ast.NewIdent("0")},
			},
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("i"),
				Op: token.LSS,
				Y:  ast.NewIdent(strconv.FormatInt(srcType.Len(), 10)),
			},
			Post: &ast.IncDecStmt{
				X:   ast.NewIdent("i"),
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{},
		}
		dstElemVar := types.NewVar(0, b.Package, dst.Name()+"[i]", dstType.Elem())
		srcElemVar := types.NewVar(0, b.Package, src.Name()+"[i]", srcType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		forStmt.Body.List = append(forStmt.Body.List, elementStmt...)
		stmts = append(stmts, forStmt)
		return stmts
	case *types.Map:
	case *types.Interface:
	case *types.Slice:
	case *types.Named:
		dstUnderType := dstType.Underlying()
		srcUnderType := src.Type().Underlying()
		var srcName = src.Name()
		switch dstUnderType.(type) {
		case *types.Basic:
			if src.Type().String() != dst.Type().String() {
				// not same type
				srcName = fmt.Sprintf("%s(%s)", importType(b.ImportName, dst), src.Name())
			}
		}
		dstUnderVar := types.NewVar(0, b.Package, dst.Name(), dstUnderType)
		srcUnderVar := types.NewVar(0, b.Package, srcName, srcUnderType)
		return b.buildStmt(dstUnderVar, srcUnderVar)
	case *types.Basic:
		srcType, ok := src.Type().(*types.Basic)
		if !ok {
			underType, ok := src.Type().Underlying().(*types.Basic)
			if !ok {
				log.Fatalf("src type is not a basic:%s", src.String())
				return nil
			}
			// cast
			castName := fmt.Sprintf("%s(%s)", importType(b.ImportName, dst), src.Name())
			src = types.NewVar(0, b.Package, castName, underType.Underlying())
			srcType = underType
		}
		if srcType.Kind() != dstType.Kind() {
			log.Fatalf("src type kind is not equal %s,%s", srcType.String(), dstType.String())
			return nil
		}
		var assignmentStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent(src.Name())},
		}
		stmts = append(stmts, assignmentStmt)
		return stmts
	default:

	}

	return stmts
}

func (b *Builder) buildPointer(dest *types.Var, src *types.Var) *ast.IfStmt {
	if _, ok := src.Type().(*types.Pointer); !ok {
		return nil
	}
	ifStmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X: &ast.Ident{
				Name: src.Name(),
			},
			Op: token.NEQ,
			Y:  &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{},
	}
	// dest = new(dest.Type)
	it := dest.Type().(*types.Pointer)
	elementVar := types.NewVar(0, b.Package, dest.Name(), it.Elem())
	initAssign := &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(dest.Name())},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{&ast.CallExpr{
			Fun:  ast.NewIdent("new"),
			Args: []ast.Expr{ast.NewIdent(importType(b.ImportName, elementVar))},
		}},
	}
	ifStmt.Body.List = append(ifStmt.Body.List, initAssign)

	elementStmt := b.buildStmt(elementVar, src)
	ifStmt.Body.List = append(ifStmt.Body.List, elementStmt...)
	return ifStmt
}

func importType(importName map[string]string, t *types.Var) string {
	typeString := t.Type().String()
	varPath := t.Pkg().Path()
	ret := strings.Replace(typeString, varPath+"/", importName[varPath], 1)
	if ret == typeString {
		// not imported
		lastIndex := strings.LastIndex(typeString, "/")
		if lastIndex == -1 {
			return typeString
		}
		ret = typeString[lastIndex+1:]
		importName[typeString[:lastIndex]] = ret
		log.Printf("import new type:%s", typeString)
	}
	return ret
}

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

//go:generate
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
	}
	builder := &Builder{
		f:        f,
		types:    srcPackage.Types,
		importer: NewImporter(),
		pkgPath:  srcPackage.PkgPath,
	}
	// parse method
	for _, method := range methods {
		builder.BuildFunc(&method)
	}
	// handle import
	builder.FillImport()
	var sb strings.Builder
	err = printer.Fprint(&sb, token.NewFileSet(), f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(sb.String())
}

func (b *Builder) BuildFunc(method *MethodAndDoc) {
	src := method.Method.Signature().Params().At(0)
	dst := method.Method.Signature().Results().At(0)
	funcName := method.AstFunc.Names[0].Name
	// add a func
	fn := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: funcName,
		},
		Type: method.AstFunc.Type.(*ast.FuncType),
		Body: &ast.BlockStmt{},
	}
	// check fn Type 参数是否有命名
	for _, field := range fn.Type.Params.List {
		if len(field.Names) == 0 {
			field.Names = []*ast.Ident{ast.NewIdent("src")}
			src = types.NewVar(0, src.Pkg(), "src", src.Type())
		}
	}
	for _, field := range fn.Type.Results.List {
		if len(field.Names) == 0 {
			field.Names = []*ast.Ident{ast.NewIdent("dst")}
			dst = types.NewVar(0, dst.Pkg(), "dst", dst.Type())
		}
	}
	b.f.Decls = append(b.f.Decls, fn)
	b.importer.ImportType(src)
	b.importer.ImportType(dst)
	stmts := b.buildStmt(dst, src)
	fn.Body.List = append(fn.Body.List, stmts...)
	fn.Body.List = append(fn.Body.List, &ast.ExprStmt{X: ast.NewIdent("return")})
}

type Builder struct {
	f        *ast.File
	types    *types.Package
	importer *Importer
	pkgPath  string
}

type Importer struct {
	pkgToName       map[string]string
	importedPkgName map[string]int
	imported        []*types.Package
}

func (i *Importer) ImportType(t *types.Var) string {
	var pkgPath string
	var pkgName string
	var typeName string
	var typPrefix string
	var pkg *types.Package
	var resolve func(tye types.Type)
	resolve = func(tye types.Type) {
		switch varType := tye.(type) {
		case *types.Named:
			pkg = varType.Obj().Pkg()
			pkgPath = pkg.Path()
			pkgName = pkg.Name()
			typeName += varType.Obj().Name()
			return
		case *types.Basic:
			typeName += varType.Name()
		case *types.Slice:
			typPrefix += "[]"
			resolve(varType.Elem())
		case *types.Pointer:
			typPrefix += "*"
			resolve(varType.Elem())
		default:
			panic("expect unreachable")
		}
	}
	resolve(t.Type())
	if pkgPath == "" {
		return typeName
	}
	pkgImportName := i.pkgToName[pkgPath]
	if pkgImportName == "" {
		pkgImportName = pkgName
		// import pkg name
		if num, ok := i.importedPkgName[pkgImportName]; ok {
			next := num + 1
			i.importedPkgName[pkgImportName] = next
			pkgImportName = pkgImportName + strconv.Itoa(next)
		} else {
			i.importedPkgName[pkgImportName] = 1
		}
		i.pkgToName[pkgPath] = pkgImportName
		i.imported = append(i.imported, pkg)
	}
	name := typPrefix + pkgName + "." + typeName
	return name
}

func NewImporter() *Importer {
	return &Importer{
		pkgToName:       map[string]string{},
		importedPkgName: map[string]int{},
	}
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
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{},
		}
		// dst = new(dst.Type)
		destPtr, srcPtr := dst.Type().(*types.Pointer), src.Type().(*types.Pointer)
		dstElemVar := types.NewVar(0, b.types, "*"+dst.Name(), destPtr.Elem())
		srcElemVar := types.NewVar(0, b.types, "*"+src.Name(), srcPtr.Elem())
		initAssign := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{
				Fun:  ast.NewIdent("new"),
				Args: []ast.Expr{ast.NewIdent(b.importer.ImportType(dstElemVar))},
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
					dstVar := types.NewVar(0, b.types, dstVarName, dstField.Type())
					srcVar := types.NewVar(0, b.types, srcVarName, srcField.Type())
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
		// for i := 0; i<n ; i++ {}
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
		dstElemVar := types.NewVar(0, b.types, dst.Name()+"[i]", dstType.Elem())
		srcElemVar := types.NewVar(0, b.types, src.Name()+"[i]", srcType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		forStmt.Body.List = append(forStmt.Body.List, elementStmt...)
		stmts = append(stmts, forStmt)
		return stmts
	case *types.Map:
		srcType, ok := src.Type().(*types.Map)
		if !ok {
			log.Fatalf("src type is not a map:%s", src.String())
			return nil
		}
		ifStmt := &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun:  ast.NewIdent("len"),
					Args: []ast.Expr{ast.NewIdent(src.Name())},
				},
				OpPos: 0,
				Op:    token.GTR,
				Y:     ast.NewIdent("0"),
			},
			Body: &ast.BlockStmt{},
		}

		dstKeyVar := types.NewVar(0, b.types, "tmpK", dstType.Key())
		dstValueVar := types.NewVar(0, b.types, "tmpV", dstType.Elem())
		srcKeyVar := types.NewVar(0, b.types, "k", srcType.Key())
		srcValueVar := types.NewVar(0, b.types, "v", srcType.Elem())

		dstKeyTypeStr := b.importer.ImportType(dstKeyVar)
		dstValueTypeStr := b.importer.ImportType(dstValueVar)

		mkStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun:    ast.NewIdent("make"),
					Lparen: 0,
					Args: []ast.Expr{
						// type
						&ast.MapType{
							Key:   ast.NewIdent(dstKeyTypeStr),
							Value: ast.NewIdent(dstValueTypeStr),
						},
						// cap
						&ast.CallExpr{
							Fun:  ast.NewIdent("len"),
							Args: []ast.Expr{ast.NewIdent(src.Name())},
						},
					},
				},
			},
		}
		ifStmt.Body.List = append(ifStmt.Body.List, mkStmt)
		// for k, v := range src.xx
		rangeStmt := &ast.RangeStmt{
			Key:   ast.NewIdent("k"),
			Value: ast.NewIdent("v"),
			Tok:   token.DEFINE,
			X:     ast.NewIdent(src.Name()),
			Body:  &ast.BlockStmt{},
		}
		ifStmt.Body.List = append(ifStmt.Body.List, rangeStmt)
		kvDeclStmt := &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent(dstKeyVar.Name())},
						Type:  ast.NewIdent(dstKeyTypeStr),
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent(dstValueVar.Name())},
						Type:  ast.NewIdent(dstValueTypeStr),
					},
				},
			}}
		// var (tmpK xx, tmpV xx)
		rangeStmt.Body.List = append(rangeStmt.Body.List, kvDeclStmt)
		assignKStmt := b.buildStmt(dstKeyVar, srcKeyVar)
		assignVStmt := b.buildStmt(dstValueVar, srcValueVar)
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignKStmt...)
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignVStmt...)
		assignMapStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{&ast.IndexExpr{
				X:      ast.NewIdent(dst.Name()),
				Lbrack: 0,
				Index:  ast.NewIdent(dstKeyVar.Name()),
				Rbrack: 0,
			}},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent(dstValueVar.Name())},
		}
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignMapStmt)
		stmts = append(stmts, ifStmt)
		return stmts
	case *types.Interface:
	case *types.Slice:
		srcType, ok := src.Type().(*types.Slice)
		if !ok {
			log.Fatalf("src type is not a array:%s", src.String())
			return nil
		}
		ifStmt := &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun:  ast.NewIdent("len"),
					Args: []ast.Expr{ast.NewIdent(src.Name())},
				},
				OpPos: 0,
				Op:    token.GTR,
				Y:     ast.NewIdent("0"),
			},
			Body: &ast.BlockStmt{},
		}
		mkStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun:    ast.NewIdent("make"),
					Lparen: 0,
					Args: []ast.Expr{
						// type
						ast.NewIdent(b.importer.ImportType(dst)),
						// cap
						&ast.CallExpr{
							Fun:  ast.NewIdent("len"),
							Args: []ast.Expr{ast.NewIdent(src.Name())},
						},
					},
					Ellipsis: 0,
					Rparen:   0,
				},
			},
		}
		ifStmt.Body.List = append(ifStmt.Body.List, mkStmt)
		// for i := 0; i<n ; i++ {}
		forStmt := &ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("i")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{ast.NewIdent("0")},
			},
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("i"),
				Op: token.LSS,
				Y: &ast.CallExpr{
					Fun:  ast.NewIdent("len"),
					Args: []ast.Expr{ast.NewIdent(src.Name())},
				},
			},
			Post: &ast.IncDecStmt{
				X:   ast.NewIdent("i"),
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{},
		}
		dstElemVar := types.NewVar(0, b.types, dst.Name()+"[i]", dstType.Elem())
		srcElemVar := types.NewVar(0, b.types, src.Name()+"[i]", srcType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		forStmt.Body.List = append(forStmt.Body.List, elementStmt...)
		ifStmt.Body.List = append(ifStmt.Body.List, forStmt)
		stmts = append(stmts, ifStmt)
		return stmts
	case *types.Named:
		dstUnderType := dstType.Underlying()
		srcUnderType := src.Type().Underlying()
		var srcName = src.Name()
		switch dstUnderType.(type) {
		case *types.Basic:
			if src.Type().String() != dst.Type().String() {
				// not same type
				srcName = fmt.Sprintf("%s(%s)", b.importer.ImportType(dst), src.Name())
			}
		}
		dstUnderVar := types.NewVar(0, b.types, dst.Name(), dstUnderType)
		srcUnderVar := types.NewVar(0, b.types, srcName, srcUnderType)
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
			castName := fmt.Sprintf("%s(%s)", b.importer.ImportType(dst), src.Name())
			src = types.NewVar(0, b.types, castName, underType.Underlying())
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

func (b *Builder) FillImport() {
	var importDecls []ast.Decl
	im := &ast.GenDecl{
		Doc:   nil,
		Tok:   token.IMPORT,
		Specs: []ast.Spec{},
	}
	for _, p := range b.importer.imported {
		spec := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "\"" + p.Path() + "\"",
			},
		}
		im.Specs = append(im.Specs, spec)
		if name, ok := b.importer.pkgToName[p.Path()]; ok && name != p.Name() {
			spec.Name = ast.NewIdent(name)
		}
	}
	importDecls = append(importDecls, im)
	b.f.Decls = append(importDecls, b.f.Decls...)
}

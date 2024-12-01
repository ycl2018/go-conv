package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Builder a file in a package
type Builder struct {
	*InitBuilder
	f        *ast.File
	types    *types.Package
	importer *Importer
	genFunc  map[string]*ast.FuncDecl
	rootNode bool
}

func NewBuilder(f *ast.File, types *types.Package, importer *Importer) *Builder {
	return &Builder{
		f:           f,
		types:       types,
		importer:    importer,
		genFunc:     make(map[string]*ast.FuncDecl),
		InitBuilder: NewInitBuilder(),
	}
}

func (b *Builder) BuildFunc(dst, src types.Type) (funcName string) {
	srcTypeName, dstTypeName := b.importer.ImportType(src), b.importer.ImportType(dst)
	funcName = b.GenFuncName(src, dst)
	b.rootNode = true
	// add a func
	fn := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: funcName,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("src")},
					Type:  ast.NewIdent(srcTypeName),
				},
			}},
			Results: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("dst")},
					Type:  ast.NewIdent(dstTypeName),
				},
			}},
		},
		Body: &ast.BlockStmt{},
	}
	b.genFunc[funcName] = fn
	srcName, dstName := "src", "dst"
	srcVar, dstVar := types.NewVar(0, b.types, srcName, src), types.NewVar(0, b.types, dstName, dst)
	stmts := b.buildStmt(dstVar, srcVar)
	fn.Body.List = append(fn.Body.List, stmts...)
	fn.Body.List = append(fn.Body.List, &ast.ExprStmt{X: ast.NewIdent("return")})
	return funcName
}

func convArrayToSlice(v types.Type) (s *types.Slice, conved, ok bool) {
	if s, ok = v.Underlying().(*types.Slice); ok {
		return s, false, true
	}
	if arr, ok := v.Underlying().(*types.Array); ok {
		return types.NewSlice(arr.Elem()), true, true
	}
	return nil, false, false
}

func convSliceToArray(v types.Type) (arr *types.Array, conved, ok bool) {
	if arr, ok = v.Underlying().(*types.Array); ok {
		return arr, false, true
	}
	if s, ok := v.Underlying().(*types.Slice); ok {
		// we can't define the length
		return types.NewArray(s.Elem(), -1), true, true
	}
	return nil, false, false
}

func convStructToPointer(v types.Type) (ptr *types.Pointer, conved, ok bool) {
	if ptr, ok := v.(*types.Pointer); ok {
		return ptr, false, true
	}
	// check if src is a Named struct
	if _, ok := v.Underlying().(*types.Struct); ok {
		return types.NewPointer(v), true, true
	}
	return nil, false, false
}

func convPtrToStruct(v types.Type) (strut *types.Struct, conved, ok bool) {
	if strut, ok := v.Underlying().(*types.Struct); ok {
		return strut, false, true
	}
	if ptr, ok := v.Underlying().(*types.Pointer); ok {
		if strut, ok := ptr.Elem().Underlying().(*types.Struct); ok {
			return strut, true, true
		}
	}
	return nil, false, false
}

func (b *Builder) buildStmt(dst *types.Var, src *types.Var) []ast.Stmt {
	defer func() {
		b.rootNode = false
	}()
	var stmts []ast.Stmt
	switch dstType := dst.Type().(type) {
	case *types.Pointer:
		srcPtrType, conved, ok := convStructToPointer(src.Type())
		if !ok {
			log.Printf("src type is not a pointer:%s", src.String())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
		}
		// check if has generated func
		if named, ok := dstType.Elem().(*types.Named); ok {
			if _, ok := named.Underlying().(*types.Struct); ok && !b.rootNode {
				funcName := b.GenFuncName(srcPtrType, dst.Type())
				convedSrcName := func() string {
					if conved {
						return "&" + src.Name()
					}
					return src.Name()
				}()
				assignStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun:  ast.NewIdent(funcName),
						Args: []ast.Expr{ast.NewIdent(convedSrcName)},
					}},
				}
				if _, ok := b.genFunc[funcName]; !ok {
					b.BuildFunc(dst.Type(), srcPtrType)
				}
				stmts = append(stmts, assignStmt)
				return stmts
			}
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
		destPtr, srcPtr := dst.Type().(*types.Pointer), srcPtrType
		srcVarName := func() string {
			if conved {
				return src.Name()
			}
			return "*" + src.Name()
		}()
		dstElemVar := types.NewVar(0, b.types, "*"+dst.Name(), destPtr.Elem())
		srcElemVar := types.NewVar(0, b.types, srcVarName, srcPtr.Elem())
		initAssign := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{
				Fun:  ast.NewIdent("new"),
				Args: []ast.Expr{ast.NewIdent(b.importer.ImportType(dstElemVar.Type()))},
			}},
		}
		ifStmt.Body.List = append(ifStmt.Body.List, initAssign)
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		ifStmt.Body.List = append(ifStmt.Body.List, elementStmt...)
		stmts = append(stmts, ifStmt)
		return stmts
	case *types.Struct:
		srcType, _, ok := convPtrToStruct(src.Type())
		if !ok {
			log.Printf("src type is not a struct:%s", src.String())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
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
				if !srcField.Exported() {
					continue
				}
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
				log.Printf("src field %s not found in struct:%s", dstFieldName, srcType.String())
				stmts = append(stmts, buildCommentExpr("omit "+dstFieldName))
			}
		}
		return stmts
	case *types.Array:
		srcArrType, _, ok := convSliceToArray(src.Type())
		if !ok {
			log.Printf("src type is not a array/slice:%s", src.String())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
		}
		// for i := 0; i<n ; i++ {}
		forStmt := &ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("i")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{ast.NewIdent("0")},
			},
			Cond: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  ast.NewIdent("i"),
					Op: token.LSS,
					Y:  ast.NewIdent(strconv.FormatInt(dstType.Len(), 10)),
				},
				Op: token.LAND,
				Y: &ast.BinaryExpr{
					X:  ast.NewIdent("i"),
					Op: token.LSS,
					Y: &ast.CallExpr{
						Fun:  ast.NewIdent("len"),
						Args: []ast.Expr{ast.NewIdent(src.Name())},
					},
				},
			},
			Post: &ast.IncDecStmt{
				X:   ast.NewIdent("i"),
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{},
		}
		dstElemVar := types.NewVar(0, b.types, dst.Name()+"[i]", dstType.Elem())
		srcElemVar := types.NewVar(0, b.types, src.Name()+"[i]", srcArrType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		forStmt.Body.List = append(forStmt.Body.List, elementStmt...)
		stmts = append(stmts, forStmt)
		return stmts
	case *types.Map:
		srcType, ok := src.Type().Underlying().(*types.Map)
		if !ok {
			log.Printf("src type is not a map:%s", src.String())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
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

		dstKeyTypeStr := b.importer.ImportType(dstKeyVar.Type())
		dstValueTypeStr := b.importer.ImportType(dstValueVar.Type())

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
		// for k, v := range src
		rangeStmt := &ast.RangeStmt{
			Key:   ast.NewIdent("k"),
			Value: ast.NewIdent("v"),
			Tok:   token.DEFINE,
			X:     ast.NewIdent(src.Name()),
			Body:  &ast.BlockStmt{},
		}
		ifStmt.Body.List = append(ifStmt.Body.List, rangeStmt)
		kDeclStmt := &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent(dstKeyVar.Name())},
						Type:  ast.NewIdent(dstKeyTypeStr),
					},
				},
			}}
		vDeclStmt := &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent(dstValueVar.Name())},
						Type:  ast.NewIdent(dstValueTypeStr),
					},
				},
			}}
		// var (tmpK xx, tmpV xx)
		rangeStmt.Body.List = append(rangeStmt.Body.List, kDeclStmt, vDeclStmt)
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
	case *types.Slice:
		srcSliceType, _, ok := convArrayToSlice(src.Type())
		if !ok {
			// check is string -> []byte/[]rune
			if db, ok := dstType.Elem().(*types.Basic); ok &&
				(db.Kind() == types.Byte || db.Kind() == types.Rune) {
				if sb, ok := src.Type().Underlying().(*types.Basic); ok && sb.Kind() == types.String {
					dstName := b.importer.ImportType(dstType)
					assignStmt := &ast.AssignStmt{
						Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{ast.NewIdent(fmt.Sprintf("%s(%s)",
							dstName,
							src.Name())),
						},
					}
					return append(stmts, assignStmt)
				}
			}
			log.Printf("src type is not a slice/array:%s", src.String())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
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
						ast.NewIdent(b.importer.ImportType(dst.Type())),
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
		srcElemVar := types.NewVar(0, b.types, src.Name()+"[i]", srcSliceType.Elem())
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
				srcName = fmt.Sprintf("%s(%s)", b.importer.ImportType(dst.Type()), src.Name())
			}
		}
		dstUnderVar := types.NewVar(0, b.types, dst.Name(), dstUnderType)
		srcUnderVar := types.NewVar(0, b.types, srcName, srcUnderType)
		return b.buildStmt(dstUnderVar, srcUnderVar)
	case *types.Basic:
		srcType, ok := src.Type().(*types.Basic)
		var valid, needCast bool
		if !ok {
			// whether dst is string, src is []byte/[]rune
			if dstType.Kind() == types.String {
				if s, ok := src.Type().(*types.Slice); ok {
					if b, ok := s.Elem().(*types.Basic); ok {
						if b.Kind() == types.Byte || b.Kind() == types.Rune {
							needCast = true
							valid = true
						}
					}
				}
			}
			if srcUnderType, ok := src.Type().Underlying().(*types.Basic); ok {
				if canCast(srcUnderType, dstType) {
					valid = true
					needCast = true
				}
			}
		} else if srcType.Kind() == dstType.Kind() { // same basic
			valid = true
		} else if canCast(srcType, dstType) { // cast basic
			valid = true
			needCast = true
		}
		if !valid {
			log.Printf("src type %s not valid", srcType.String())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
		}
		var srcName = src.Name()
		// check if already cast in parent types.Named var
		if needCast && strings.Index(src.Name(), "(") == -1 {
			srcName = fmt.Sprintf("%s(%s)", b.importer.ImportType(dst.Type()), src.Name())
		}
		var assignmentStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent(srcName)},
		}
		stmts = append(stmts, assignmentStmt)
		return stmts
	default:
		stmts = append(stmts, buildCommentExpr("omit "+dst.Name()))
	}

	return stmts
}

func canCast(from, to *types.Basic) bool {
	fromInfo, toInfo := from.Info(), to.Info()
	// numbers
	if (fromInfo|types.IsNumeric|types.IsUnsigned) != 0 && (toInfo|types.IsNumeric|types.IsUnsigned) != 0 {
		return true
	}
	return false
}

func (b *Builder) Generate() ([]byte, error) {
	// fill and sort import
	b.fillImport()
	// sort func by name
	var funcNames []string
	for name := range b.genFunc {
		funcNames = append(funcNames, name)
	}
	sort.Strings(funcNames)
	for _, name := range funcNames {
		b.f.Decls = append(b.f.Decls, b.genFunc[name])
	}
	// add init func
	initFunc := b.GenInit()
	b.f.Decls = append(b.f.Decls, initFunc)

	var sb bytes.Buffer
	sb.WriteString("// Code generated by github.com/ycl2018/go-conv DO NOT EDIT.\n\n")
	// format
	err := format.Node(&sb, token.NewFileSet(), b.f)
	if err != nil {
		return nil, fmt.Errorf("[go-conv]: failed to format code: %w", err)
	}
	return sb.Bytes(), nil
}

func (b *Builder) fillImport() {
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

func buildCommentExpr(comment string) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "// " + comment,
		},
	}
}

func cleanName(name string) string {
	var s strings.Builder
	var first = true
	for _, c := range name {
		if unicode.IsLetter(c) || (unicode.IsNumber(c) && !first) {
			if first {
				s.WriteString(strings.ToUpper(string(c)))
				first = false
			} else {
				s.WriteRune(c)
			}
		}
	}
	return s.String()
}

func (b *Builder) GenFuncName(src, dst types.Type) string {
	srcTypeName, dstTypeName := b.importer.ImportType(src), b.importer.ImportType(dst)
	return cleanName(srcTypeName) + "To" + cleanName(dstTypeName)
}

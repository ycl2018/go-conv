package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type BuildMode int

const (
	BuildModeCopy BuildMode = iota + 1
	BuildModeConv
)

func (m BuildMode) String() string {
	switch m {
	case BuildModeCopy:
		return "copyMode"
	case BuildModeConv:
		return "convMode"
	}
	return "_"
}

// Builder a file in a package
type Builder struct {
	*InitFuncBuilder
	f           *ast.File
	types       *types.Package
	importer    *Importer
	genFunc     map[string]*ast.FuncDecl
	rootNode    bool
	buildConfig BuildConfig
	logger      *Logger
}

func NewBuilder(f *ast.File, types *types.Package, logger *Logger) *Builder {
	return &Builder{
		f:               f,
		types:           types,
		importer:        NewImporter(types.Path()),
		genFunc:         make(map[string]*ast.FuncDecl),
		InitFuncBuilder: NewInitFuncBuilder(),
		logger:          logger,
	}
}

func (b *Builder) BuildFunc(dst, src types.Type, buildConfig BuildConfig) (funcName string) {
	srcTypeName, dstTypeName := b.importer.ImportType(src), b.importer.ImportType(dst)
	funcName = b.GenFuncName(src, dst, buildConfig)
	b.logger.Printf("generate function:%s by %s", funcName, buildConfig)
	b.rootNode = true
	b.buildConfig = buildConfig
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

func isStruct(v types.Type) (named, ok bool) {
	if _, ok := v.(*types.Struct); ok {
		return false, true
	}
	// check if src is a Named struct
	if namedTypes, ok := v.(*types.Named); ok {
		if _, ok := namedTypes.Underlying().(*types.Struct); ok {
			return true, true

		}
	}
	return false, false
}

func convPtrToStruct(v types.Type) (strut *types.Struct, isPtr, ok bool) {
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

func (b *Builder) _shallowCopy(dst, src *types.Var) ([]ast.Stmt, bool) {
	// exactly same type
	var stmts []ast.Stmt
	if types.AssignableTo(src.Type(), dst.Type()) {
		var assignmentStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent(src.Name())},
		}
		stmts = append(stmts, assignmentStmt)
		return stmts, true
	}
	if types.ConvertibleTo(src.Type(), dst.Type()) {
		convertName := fmt.Sprintf("(%s)(%s)", b.importer.ImportType(dst.Type()), src.Name())
		assignStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent(convertName)},
		}
		return append(stmts, assignStmt), true
	}
	return nil, false
}

func (b *Builder) buildStmt(dst *types.Var, src *types.Var) []ast.Stmt {
	defer func() {
		b.rootNode = false
	}()
	var stmts []ast.Stmt
	if b.buildConfig.BuildMode == BuildModeConv {
		if ret, ok := b._shallowCopy(dst, src); ok {
			return ret
		}
	}
	switch dstType := dst.Type().(type) {
	case *types.Pointer:
		elemType, ptrDepth, srcIsPtr := dePointer(src.Type())
		_, srcIsStruct := isStruct(elemType)
		// check has generated func
		if !b.rootNode && srcIsStruct && isPointerToStruct(dstType) {
			funcName := b.GenFuncName(elemType, dst.Type(), b.buildConfig)
			convedSrcName := func() string {
				if !srcIsPtr {
					return addressName(src.Name(), 1)
				} else {
					ptrToName(src.Name(), ptrDepth-1)
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
				b.BuildFunc(dst.Type(), types.NewPointer(elemType), b.buildConfig)
			}
			stmts = append(stmts, assignStmt)
			return stmts
		}

		_, depth, _ := dePointer(dstType)
		if depth != 1 {
			b.logger.Printf("omit %s :only support one level pointer", dst.Name())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
		}
		dstName := ptrToName(dst.Name(), depth)
		srcName := ptrToName(src.Name(), ptrDepth)
		dstElemVar := types.NewVar(0, b.types, dstName, dstType.Elem())
		srcElemVar := types.NewVar(0, b.types, srcName, elemType)
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		// dst = new(dst.Type)
		initAssign := &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{
				Fun:  ast.NewIdent("new"),
				Args: []ast.Expr{ast.NewIdent(b.importer.ImportType(dstElemVar.Type()))},
			}},
		}
		if srcIsPtr {
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
			ifStmt.Body.List = append(ifStmt.Body.List, initAssign)
			ifStmt.Body.List = append(ifStmt.Body.List, elementStmt...)
			stmts = append(stmts, ifStmt)
		} else {
			stmts = append(stmts, initAssign)
			stmts = append(stmts, elementStmt...)
		}
		return stmts
	case *types.Named:
		dstUnderType := dstType.Underlying()
		var srcType = src.Type()
		var srcName = src.Name()
		switch dstUnderType.(type) {
		case *types.Basic: // for named basic
			srcElemType, ptrDepth, isPtr := dePointer(src.Type())
			if !isPtr && types.ConvertibleTo(srcType, dstType) { // not a Pointer
				srcName = fmt.Sprintf("%s(%s)", b.importer.ImportType(dst.Type()), src.Name())
				assignStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{ast.NewIdent(srcName)},
				}
				return append(stmts, assignStmt)
			}
			if isPtr && ptrDepth == 1 && types.ConvertibleTo(srcElemType, dstType) { // a Pointer
				ifStmt := &ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent(src.Name()),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{},
				}
				assignStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{ast.NewIdent(ptrToName(srcName, ptrDepth))},
				}
				ifStmt.Body.List = append(ifStmt.Body.List, assignStmt)
				return append(stmts, assignStmt)
			}
			b.logger.Printf("omit %s :named type can't cast from %s", dst.Name(), src.Name())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
		}
		dstUnderVar := types.NewVar(0, b.types, dst.Name(), dstUnderType)
		srcUnderVar := types.NewVar(0, b.types, srcName, srcType)
		return b.buildStmt(dstUnderVar, srcUnderVar)
	case *types.Struct:
		srcType, isPtr, ok := convPtrToStruct(src.Type())
		if !ok {
			b.logger.Printf("omit %s :%s type is not a struct/pointer to struct", dst.Name(), src.Name())
			return append(stmts, buildCommentExpr("omit "+dst.Name()))
		}
		// right side a struct Pointer, de pointer to
		if isPtr {
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
			ptrToSrcName := ptrToName(src.Name(), 1)
			dstElemVar := types.NewVar(0, b.types, dst.Name(), dstType)
			srcElemVar := types.NewVar(0, b.types, ptrToSrcName, srcType)
			elementStmt := b.buildStmt(dstElemVar, srcElemVar)
			ifStmt.Body.List = append(ifStmt.Body.List, elementStmt...)
			stmts = append(stmts, ifStmt)
			return stmts
		}
		srcName := strings.TrimPrefix(src.Name(), "*")
		dstName := strings.TrimPrefix(dst.Name(), "*") // for struct type, compiler can de Pointer automatically
		for i := range dstType.NumFields() {
			dstField := dstType.Field(i)
			if !dstField.Exported() {
				continue
			}
			dstFieldName := dstField.Name()
			if dstField.Embedded() {
				dstVar := types.NewVar(0, b.types, dstName+"."+dstFieldName, dstField.Type())
				srcVar := types.NewVar(0, b.types, src.Name(), src.Type())
				fieldStmt := b.buildStmt(dstVar, srcVar)
				stmts = append(stmts, fieldStmt...)
				continue
			}
			// match srcField
			if srcField, ok := matchField(dstField, srcType); ok {
				dstVarName := dstName + "." + dstFieldName
				srcVarName := srcName + "." + srcField.Name()
				dstVar := types.NewVar(0, b.types, dstVarName, dstField.Type())
				srcVar := types.NewVar(0, b.types, srcVarName, srcField.Type())
				fieldStmt := b.buildStmt(dstVar, srcVar)
				stmts = append(stmts, fieldStmt...)
			} else {
				b.logger.Printf("omit %s :not find match field in %s", dstFieldName, srcName)
				stmts = append(stmts, buildCommentExpr("omit "+dstFieldName))
			}
		}
		return stmts
	case *types.Array:
		srcArrType, _, ok := convSliceToArray(src.Type())
		if !ok {
			b.logger.Printf("omit %s :%s type is not a array/slice", dst.Name(), src.Name())
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
			b.logger.Printf("omit %s :%s type is not a map", dst.Name(), src.Name())
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
			b.logger.Printf("omit %s :%s type is not a slice/array", dst.Name(), src.Name())
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
	case *types.Basic:
		if types.AssignableTo(src.Type(), dst.Type()) {
			var assignmentStmt = &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{ast.NewIdent(src.Name())},
			}
			stmts = append(stmts, assignmentStmt)
			return stmts
		}
		if types.ConvertibleTo(src.Type(), dst.Type()) {
			convertName := fmt.Sprintf("%s(%s)", b.importer.ImportType(dst.Type()), src.Name())
			assignStmt := &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{ast.NewIdent(convertName)},
			}
			return append(stmts, assignStmt)
		}
		// check if src pointer to elem can convert to dst
		srcElemType, ptrDepth, srcIsPtr := dePointer(src.Type())

		if srcIsPtr && ptrDepth == 1 && types.ConvertibleTo(srcElemType, dstType) { // a Pointer
			ifStmt := &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  ast.NewIdent(src.Name()),
					Op: token.NEQ,
					Y:  ast.NewIdent("nil"),
				},
				Body: &ast.BlockStmt{},
			}
			srcPtrToName := ptrToName(src.Name(), ptrDepth)
			if !types.AssignableTo(srcElemType, dstType) {
				srcPtrToName = fmt.Sprintf("%s(%s)", dstType.String(), srcPtrToName)
			}
			assignStmt := &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent(dst.Name())},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{ast.NewIdent(srcPtrToName)},
			}
			ifStmt.Body.List = append(ifStmt.Body.List, assignStmt)
			return append(stmts, assignStmt)
		}
		b.logger.Printf("omit %s :basic type can't cast from %s (or it pointers to)", dst.Name(), src.Name())
		return append(stmts, buildCommentExpr("omit "+dst.Name()))
	default:
		b.logger.Printf("omit %s :type not support yet", dst.Name())
		stmts = append(stmts, buildCommentExpr("omit "+dst.Name()))
	}

	return stmts
}

func isPointerToStruct(p *types.Pointer) bool {
	if elem := p.Elem(); elem != nil {
		if _, ok := elem.Underlying().(*types.Struct); ok {
			return ok
		}
	}
	return false
}

func dePointer(t types.Type) (elemType types.Type, ptrDepth int, isPtr bool) {
	var ret = t
	ptr, _ := t.(*types.Pointer)
	for ptr != nil {
		ptrDepth++
		ret = ptr.Elem()
		v, ok := ptr.Elem().(*types.Pointer)
		if ok {
			ptr = v
		} else {
			ptr = nil
		}
	}

	return ret, ptrDepth, ptrDepth > 0
}

func ptrToName(name string, ptrDepth int) string {
	for ptrDepth > 0 {
		if strings.HasPrefix(name, "&") {
			name = strings.TrimPrefix(name, "&")
		} else {
			name = "*" + name
		}
		ptrDepth--
	}
	return name
}

func addressName(name string, addrDepth int) string {
	for addrDepth > 0 {
		if strings.HasPrefix(name, "*") {
			name = strings.TrimPrefix(name, "*")
		} else {
			name = "&" + name
		}
		addrDepth--
	}
	return name
}

// matchField find a matched Field in srcStruct with dstField
func matchField(dstField *types.Var, srcStruct *types.Struct) (matched *types.Var, match bool) {
	// by name
	for i := range srcStruct.NumFields() {
		srcField := srcStruct.Field(i)
		if !srcField.Exported() {
			continue
		}
		if srcField.Name() == dstField.Name() {
			return srcField, true
		}
		if srcField.Embedded() {
			if embedStruct, ok := srcField.Type().Underlying().(*types.Struct); ok {
				if v, ok := matchField(dstField, embedStruct); ok {
					return v, true
				}
			}
		}
	}
	return nil, false
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
	if len(im.Specs) > 0 {
		importDecls = append(importDecls, im)
	}
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

func (b *Builder) GenFuncName(src, dst types.Type, buildConfig BuildConfig) string {
	srcTypeName, dstTypeName := b.importer.ImportType(src), b.importer.ImportType(dst)
	switch buildConfig.BuildMode {
	case BuildModeConv:
		// default
	case BuildModeCopy:
		return "Copy" + cleanName(srcTypeName) + "To" + cleanName(dstTypeName)
	}
	return cleanName(srcTypeName) + "To" + cleanName(dstTypeName)
}

package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"
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
	fieldPath   path
	logger      *Logger
}

func NewBuilder(f *ast.File, types *types.Package) *Builder {
	return &Builder{
		f:               f,
		types:           types,
		importer:        NewImporter(types.Path()),
		genFunc:         make(map[string]*ast.FuncDecl),
		InitFuncBuilder: NewInitFuncBuilder(),
		logger:          DefaultLogger,
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

func (b *Builder) _shallowCopy(dst, src *types.Var) ([]ast.Stmt, bool) {
	var stmts []ast.Stmt
	// exactly same type
	if src.Type().String() == dst.Type().String() {
		var assignmentStmt = buildAssignStmt(dst.Name(), src.Name())
		stmts = append(stmts, assignmentStmt)
		return stmts, true
	}
	if types.ConvertibleTo(src.Type(), dst.Type()) {
		convertName := fmt.Sprintf("%s(%s)", parenthesesName(b.importer.ImportType(dst.Type())), src.Name())
		assignStmt := buildAssignStmt(dst.Name(), convertName)
		return append(stmts, assignStmt), true
	}
	return nil, false
}

func (b *Builder) buildStmt(dst *types.Var, src *types.Var) []ast.Stmt {
	defer func() {
		b.rootNode = false
	}()
	var stmts []ast.Stmt
	// ignore
	for _, ignoreType := range b.buildConfig.Ignore {
		if b.fieldPath.matchIgnore(ignoreType, src.Type()) {
			b.logger.Printf("apply ignore on %s", src.Name())
			b.buildCommentExpr(&stmts, "apply ignore option on %s", src.Name())
			return stmts
		}
	}
	// transfer
	for _, transfer := range b.buildConfig.Transfer {
		if b.fieldPath.matchTransfer(transfer, dst, src) {
			b.logger.Printf("apply transfer on %s", src.Name())
			b.buildCommentExpr(&stmts, "apply transfer option on %s", transfer.FuncName)
			assignStmt := buildAssignStmt(dst.Name(), fmt.Sprintf("%s(%s)", transfer.FuncName, src.Name()))
			stmts = append(stmts, assignStmt)
			return stmts
		}
	}
	// filter
	for _, filter := range b.buildConfig.Filter {
		if b.fieldPath.matchFilter(filter, src.Type()) {
			b.logger.Printf("apply filter on %s", src.Name())
			b.buildCommentExpr(&stmts, "apply filter option on %s", filter.FuncName)
			newSrcName := "filtered" + cleanName(src.Name())
			assignStmt := buildDefineStmt(newSrcName, fmt.Sprintf("%s(%s)", filter.FuncName, src.Name()))
			src = types.NewVar(0, b.types, newSrcName, src.Type())
			stmts = append(stmts, assignStmt)
		}
	}

	if b.buildConfig.BuildMode == BuildModeConv {
		if ret, ok := b._shallowCopy(dst, src); ok {
			return append(stmts, ret...)
		}
	}
	switch dstType := dst.Type().(type) {
	case *types.Pointer:
		elemType, ptrDepth, srcIsPtr := dePointer(src.Type())
		_, srcIsStruct := isStruct(elemType)
		// check has generated func
		if !b.rootNode && srcIsStruct && isPointerToStruct(dstType) {
			funcName := b.GenFuncName(types.NewPointer(elemType), dst.Type(), b.buildConfig)
			convSrcName := func() string {
				if !srcIsPtr {
					return addressName(src.Name(), 1)
				} else {
					ptrToName(src.Name(), ptrDepth-1)
				}
				return src.Name()
			}()
			assignStmt := buildAssignStmt(dst.Name(), fmt.Sprintf("%s(%s)", funcName, convSrcName))
			if _, ok := b.genFunc[funcName]; !ok {
				b.BuildFunc(dst.Type(), types.NewPointer(elemType), b.buildConfig)
			}
			stmts = append(stmts, assignStmt)
			return stmts
		}
		_, depth, _ := dePointer(dstType)
		if depth != 1 {
			b.logger.Printf("omit %s :only support one level pointer", dst.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}
		dstName := ptrToName(dst.Name(), depth)
		var srcName = src.Name()
		if !srcIsStruct {
			srcName = parenthesesName(ptrToName(src.Name(), ptrDepth))
		}
		dstElemVar := types.NewVar(0, b.types, dstName, dstType.Elem())
		srcElemVar := types.NewVar(0, b.types, srcName, elemType)
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		// dst = new(dst.Type)
		initAssign := buildAssignStmt(dst.Name(), fmt.Sprintf("new(%s)", b.importer.ImportType(dstElemVar.Type())))
		if srcIsPtr {
			ifStmt := buildIfStmt(src.Name(), token.NEQ, "nil")
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
		switch dstUnderType.(type) {
		case *types.Basic, *types.Struct: // for named basic/struct type is a special basic type
			if _, ok := dstUnderType.(*types.Struct); ok {
				if b.buildConfig.BuildMode != BuildModeConv {
					break
				}
			}
			srcElemType, ptrDepth, isPtr := dePointer(src.Type())
			var srcName = src.Name()
			if ptrDepth < 2 && types.ConvertibleTo(srcElemType, dstType) {
				b.logger.Printf("convert:%s to %s", srcElemType.String(), dstType.String())
				srcName = fmt.Sprintf("%s(%s)", parenthesesName(b.importer.ImportType(dst.Type())), ptrToName(srcName, ptrDepth))
				if !isPtr { // not a Pointer
					assignStmt := buildAssignStmt(dst.Name(), srcName)
					return append(stmts, assignStmt)
				} else {
					ifStmt := buildIfStmt(src.Name(), token.NEQ, "nil")
					assignStmt := buildAssignStmt(dst.Name(), srcName)
					ifStmt.Body.List = append(ifStmt.Body.List, assignStmt)
					return append(stmts, assignStmt)
				}
			}
		}
		dstUnderVar := types.NewVar(0, b.types, dst.Name(), dstUnderType)
		return b.buildStmt(dstUnderVar, src)
	case *types.Struct:
		srcStructType, isPtr, ok := convPtrToStruct(src.Type())
		if !ok {
			b.logger.Printf("omit %s :%s type is not a struct/pointer to struct", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}
		if isPtr {
			stmts = append(stmts, b.dePointerSrcStmt(dst, src, srcStructType)...)
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
				fieldStmt := b.buildStmt(dstVar, src)
				stmts = append(stmts, fieldStmt...)
				continue
			}
			// match srcField
			if srcField, ok := b.matchField(dstField, srcStructType, src.Type().String()); ok {
				dstVarName := dstName + "." + dstFieldName
				srcVarName := srcName + "." + srcField.Name()
				b.logger.Printf("try assign [%s -> %s]\ttype [%s -> %s]", srcVarName, dstVarName,
					srcField.Type().String(), dstField.Type().String())
				dstVar := types.NewVar(0, b.types, dstVarName, dstField.Type())
				srcVar := types.NewVar(0, b.types, srcVarName, srcField.Type())
				fieldStmt := b.buildStmt(dstVar, srcVar)
				stmts = append(stmts, fieldStmt...)
				b.fieldPath.Pop()
			} else {
				b.logger.Printf("omit %s :not find match field in %s", dstFieldName, srcName)
				b.buildCommentExpr(&stmts, "omit "+dstFieldName)
			}
		}
		return stmts
	case *types.Array:
		srcElemType, ptrDepth, isPtr := dePointer(src.Type())
		srcArrType, isSlice, ok := convSliceToArray(srcElemType)
		if !ok || ptrDepth > 1 {
			b.logger.Printf("omit %s :%s type is not a array/slice", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}
		if isPtr {
			return b.dePointerSrcStmt(dst, src, srcElemType)
		}
		var srcName = src.Name()
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
				Y:  ast.NewIdent(strconv.FormatInt(dstType.Len(), 10)),
			},
			Post: &ast.IncDecStmt{
				X:   ast.NewIdent("i"),
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{},
		}
		if isSlice || srcArrType.Len() > dstType.Len() {
			forStmt.Cond = &ast.BinaryExpr{
				X:  forStmt.Cond,
				Op: token.LAND,
				Y: &ast.BinaryExpr{
					X:  ast.NewIdent("i"),
					Op: token.LSS,
					Y: &ast.CallExpr{
						Fun:  ast.NewIdent("len"),
						Args: []ast.Expr{ast.NewIdent(srcName)},
					},
				},
			}
		}
		dstElemVar := types.NewVar(0, b.types, parenthesesName(dst.Name())+"[i]", dstType.Elem())
		srcElemVar := types.NewVar(0, b.types, srcName+"[i]", srcArrType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		forStmt.Body.List = append(forStmt.Body.List, elementStmt...)
		stmts = append(stmts, forStmt)
		return stmts
	case *types.Map:
		srcElemType, ptrDepth, isPtr := dePointer(src.Type())
		srcType, ok := srcElemType.Underlying().(*types.Map)
		if !ok || ptrDepth > 1 {
			b.logger.Printf("omit %s :%s type is not a map", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}
		if isPtr {
			return b.dePointerSrcStmt(dst, src, srcElemType)
		}
		var srcName = parenthesesName(ptrToName(src.Name(), ptrDepth))
		ifStmt := buildIfStmt(fmt.Sprintf("len(%s)", srcName), token.GTR, "0")
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
							Args: []ast.Expr{ast.NewIdent(srcName)},
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
			X:     ast.NewIdent(srcName),
			Body:  &ast.BlockStmt{},
		}
		ifStmt.Body.List = append(ifStmt.Body.List, rangeStmt)
		kDeclStmt := buildVarDecl(dstKeyVar.Name(), dstKeyTypeStr)
		vDeclStmt := buildVarDecl(dstValueVar.Name(), dstValueTypeStr)
		// var (tmpK xx, tmpV xx)
		rangeStmt.Body.List = append(rangeStmt.Body.List, kDeclStmt, vDeclStmt)
		assignKStmt := b.buildStmt(dstKeyVar, srcKeyVar)
		assignVStmt := b.buildStmt(dstValueVar, srcValueVar)
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignKStmt...)
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignVStmt...)
		assignMapStmt := buildAssignStmt(fmt.Sprintf("%s[%s]", parenthesesName(dst.Name()), dstKeyVar.Name()), dstValueVar.Name())
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignMapStmt)
		stmts = append(stmts, ifStmt)
		return stmts
	case *types.Slice:
		srcElemType, ptrDepth, isPtr := dePointer(src.Type())
		srcSliceType, isArray, ok := convArrayToSlice(srcElemType)
		if !ok || ptrDepth > 1 {
			// check is string -> []byte/[]rune
			if db, ok := dstType.Elem().(*types.Basic); ok &&
				(db.Kind() == types.Byte || db.Kind() == types.Rune) {
				if sb, ok := src.Type().Underlying().(*types.Basic); ok && sb.Kind() == types.String {
					dstName := b.importer.ImportType(dstType)
					assignStmt := buildAssignStmt(dst.Name(), fmt.Sprintf("%s(%s)", dstName,
						src.Name()))
					return append(stmts, assignStmt)
				}
			}
			b.logger.Printf("omit %s :%s type is not a slice/array", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}
		if isPtr {
			return b.dePointerSrcStmt(dst, src, srcElemType)
		}
		var srcName = src.Name()
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
							Args: []ast.Expr{ast.NewIdent(srcName)},
						},
					},
					Ellipsis: 0,
					Rparen:   0,
				},
			},
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
				Y: &ast.CallExpr{
					Fun:  ast.NewIdent("len"),
					Args: []ast.Expr{ast.NewIdent(srcName)},
				},
			},
			Post: &ast.IncDecStmt{
				X:   ast.NewIdent("i"),
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{},
		}
		dstElemVar := types.NewVar(0, b.types, parenthesesName(dst.Name())+"[i]", dstType.Elem())
		srcElemVar := types.NewVar(0, b.types, srcName+"[i]", srcSliceType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		forStmt.Body.List = append(forStmt.Body.List, elementStmt...)
		if !isArray {
			ifStmt := buildIfStmt(fmt.Sprintf("len(%s)", srcName), token.GTR, "0")
			stmts = append(stmts, ifStmt)
			ifStmt.Body.List = append(ifStmt.Body.List, mkStmt)
			ifStmt.Body.List = append(ifStmt.Body.List, forStmt)
		} else {
			stmts = append(stmts, mkStmt)
			stmts = append(stmts, forStmt)
		}
		return stmts
	case *types.Basic:
		// check if src pointer to elem can convert to dst
		srcElemType, ptrDepth, srcIsPtr := dePointer(src.Type())
		if ptrDepth < 2 && types.ConvertibleTo(srcElemType, dstType) {
			if srcIsPtr {
				return b.dePointerSrcStmt(dst, src, srcElemType)
			}
			srcName := src.Name()
			if !types.AssignableTo(srcElemType, dstType) {
				// need cast
				srcName = fmt.Sprintf("%s(%s)", parenthesesName(b.importer.ImportType(dst.Type())), srcName)
			}
			assignStmt := buildAssignStmt(dst.Name(), srcName)
			stmts = append(stmts, assignStmt)
			return stmts
		}
		b.logger.Printf("omit %s :basic type can't cast from %s (or it pointers to)", dst.Name(), src.Name())
		b.buildCommentExpr(&stmts, "omit "+dst.Name())
	default:
		b.logger.Printf("omit %s :type not support yet", dst.Name())
		b.buildCommentExpr(&stmts, "omit "+dst.Name())
	}

	return stmts
}

func (b *Builder) buildCommentExpr(stmts *[]ast.Stmt, format string, args ...any) {
	if b.buildConfig.NoComment {
		return
	}
	*stmts = append(*stmts, &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "// " + fmt.Sprintf(format, args...),
		}})
}

func (b *Builder) dePointerSrcStmt(dst *types.Var, src *types.Var, srcElemType types.Type) []ast.Stmt {
	ifStmt := buildIfStmt(src.Name(), token.NEQ, "nil")
	var needParentheses = true
	switch srcElemType.Underlying().(type) {
	case *types.Struct:
		needParentheses = false
	case *types.Basic:
		needParentheses = false
	}
	ptrToSrcName := ptrToName(src.Name(), 1)
	if needParentheses {
		ptrToSrcName = parenthesesName(ptrToSrcName)
	}
	srcElemVar := types.NewVar(0, b.types, ptrToSrcName, srcElemType)
	elementStmt := b.buildStmt(dst, srcElemVar)
	ifStmt.Body.List = append(ifStmt.Body.List, elementStmt...)
	return []ast.Stmt{ifStmt}
}

func parenthesesName(name string) string {
	if strings.HasPrefix(name, "*") {
		return "(" + name + ")"
	}
	return name
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

func convArrayToSlice(v types.Type) (s *types.Slice, isArray, ok bool) {
	if s, ok = v.Underlying().(*types.Slice); ok {
		return s, false, true
	}
	if arr, ok := v.Underlying().(*types.Array); ok {
		return types.NewSlice(arr.Elem()), true, true
	}
	return nil, false, false
}

func convSliceToArray(v types.Type) (arr *types.Array, isSlice, ok bool) {
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

// matchField find a matched Field in srcStruct with dstField
func (b *Builder) matchField(dstField *types.Var, srcStruct *types.Struct, srcTypeString string) (
	matched *types.Var, match bool) {
	// by name
	for i := range srcStruct.NumFields() {
		srcField := srcStruct.Field(i)
		if !srcField.Exported() {
			continue
		}
		matchFromField := srcField.Name()
		if setMatch, ok := b.buildConfig.FieldMatcher.HasMatch(srcTypeString, matchFromField); ok {
			matchFromField = setMatch
		}
		if matchFromField == dstField.Name() || (b.buildConfig.CaseInsensitive &&
			(strings.ToUpper(matchFromField) == strings.ToUpper(dstField.Name()))) {
			b.fieldPath.Push(fieldStep{name: srcField.Name(), structName: srcTypeString})
			return srcField, true
		}
		if srcField.Embedded() {
			if embedStruct, ok := srcField.Type().Underlying().(*types.Struct); ok {
				if v, ok := b.matchField(dstField, embedStruct, srcField.Type().String()); ok {
					return v, true
				}
			}
		}
	}
	return nil, false
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
	if !b.buildConfig.NoInit {
		initFunc := b.GenInit()
		b.f.Decls = append(b.f.Decls, initFunc)
	}

	// format
	var buf bytes.Buffer
	fileSet := token.NewFileSet()
	err := printer.Fprint(&buf, fileSet, b.f)
	if err != nil {
		return nil, fmt.Errorf("format.Node internal error (%s)", err)
	}
	// parse
	const parserMode = parser.ParseComments | parser.SkipObjectResolution
	file, err := parser.ParseFile(fileSet, "", buf.Bytes(), parserMode)
	if err != nil {
		// We should never get here. If we do, provide good diagnostic.
		return nil, fmt.Errorf("format.Node internal error (%s)", err)
	}
	ast.SortImports(fileSet, file)
	var sb bytes.Buffer
	sb.WriteString("// Code generated by github.com/ycl2018/go-conv DO NOT EDIT.\n\n")
	err = printer.Fprint(&sb, fileSet, file)
	if err != nil {
		return nil, fmt.Errorf("format.Node internal error (%s)", err)
	}
	return sb.Bytes(), nil
}

func (b *Builder) fillImport() {
	var importDecls = b.importer.GenImportDecl()
	b.f.Decls = append(importDecls, b.f.Decls...)
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

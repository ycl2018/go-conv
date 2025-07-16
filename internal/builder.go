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
	genFunc     map[string]*BuildFunc
	rootNode    bool
	curGenFunc  string
	scope       *Scope
	buildConfig BuildConfig // buildConfig differs from var to var
	fieldPath   path
	varToNamed  map[*types.Var]*types.Named
	logger      *Logger
}

func NewBuilder(f *ast.File, pkg *types.Package) *Builder {
	return &Builder{
		f:               f,
		types:           pkg,
		importer:        NewImporter(pkg.Path()),
		genFunc:         make(map[string]*BuildFunc),
		InitFuncBuilder: NewInitFuncBuilder(),
		logger:          DefaultLogger,
		scope:           NewScope("global", nil),
		varToNamed:      make(map[*types.Var]*types.Named),
	}
}

func (b *Builder) pushScope(scopeName string) {
	b.scope = NewScope(scopeName, b.scope)
	b.logger.Printf("push scope:%s", scopeName)
}

func (b *Builder) popScope() {
	b.logger.Printf("pop scope:%s", b.scope.Name)
	b.scope = b.scope.EnclosingScope
}

func (b *Builder) newVar(name string, tye types.Type) *types.Var {
	return types.NewVar(0, b.types, name, tye)
}

func (b *Builder) BuildFunc(dst, src types.Type, buildConfig BuildConfig) (funcName string) {
	srcTypeName, dstTypeName := b.importer.ImportType(src), b.importer.ImportType(dst)
	funcName = b.GenFuncName(src, dst, buildConfig)
	b.logger.Printf("generate function:%s by %s", funcName, buildConfig)
	b.curGenFunc = funcName
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
	b.genFunc[funcName] = &BuildFunc{
		GenFunc:     fn,
		buildConfig: &buildConfig,
	}
	// funcScope's enclosing scope is global scope
	funcScope := NewScope("func@"+funcName, nil)
	prevScope := b.scope
	b.scope = funcScope
	srcName, dstName := "src", "dst"
	srcVar, dstVar := b.newVar(srcName, src), b.newVar(dstName, dst)
	stmts := b.buildStmt(dstVar, srcVar)
	b.scope = prevScope
	fn.Body.List = append(fn.Body.List, stmts...)
	fn.Body.List = append(fn.Body.List, &ast.ExprStmt{X: ast.NewIdent("return")})
	return funcName
}

func (b *Builder) _shallowCopy(dst, src *types.Var) ([]ast.Stmt, bool) {
	var stmts []ast.Stmt
	// struct can convert to struct directly
	// we need check if the struct has match ignore fields
	checkIgnore := func() (matched bool) {
		elemType, _, _ := dePointer(src.Type())
		_, ok := isStruct(elemType)
		if !ok {
			return false
		}
		for _, ignoreType := range b.buildConfig.Ignore {
			if len(ignoreType.Fields) > 0 &&
				b.fieldPath.matchIgnore(IgnoreType{
					Tye:        ignoreType.Tye,
					Paths:      ignoreType.Paths,
					IgnoreSide: ignoreType.IgnoreSide,
				}, elemType, dst.Type()) {
				return true
			}
		}
		return false
	}
	// exactly same type
	if src.Type().String() == dst.Type().String() && !checkIgnore() {
		var assignmentStmt = buildAssignStmt(dst.Name(), src.Name())
		stmts = append(stmts, assignmentStmt)
		return stmts, true
	}
	if types.ConvertibleTo(src.Type(), dst.Type()) && !checkIgnore() {
		convertName := fmt.Sprintf("%s(%s)", parenthesesName(b.importer.ImportType(dst.Type())), src.Name())
		assignStmt := buildAssignStmt(dst.Name(), convertName)
		return append(stmts, assignStmt), true
	}
	return nil, false
}

func (b *Builder) buildStmt(dst *types.Var, src *types.Var) []ast.Stmt {
	var stmts []ast.Stmt
	if _, ok := dst.Type().(*types.Pointer); !ok {
		b.rootNode = false
	}
	// ignore
	for _, ignoreType := range b.buildConfig.Ignore {
		if b.fieldPath.matchIgnore(ignoreType, src.Type(), dst.Type()) {
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
			newSrcName := "transferred" + cleanName(src.Name())
			assignStmt := buildDefineStmt(newSrcName, fmt.Sprintf("%s(%s)", transfer.FuncName, src.Name()))
			src = b.newVar(newSrcName, dst.Type())
			stmts = append(stmts, assignStmt)
		}
	}
	// filter
	for _, filter := range b.buildConfig.Filter {
		if b.fieldPath.matchFilter(filter, src.Type(), SideSrc) {
			b.logger.Printf("apply filter on %s", src.Name())
			b.buildCommentExpr(&stmts, "apply filter option on %s", filter.FuncName)
			newSrcName := "filtered" + cleanName(src.Name())
			assignStmt := buildDefineStmt(newSrcName, fmt.Sprintf("%s(%s)", filter.FuncName, src.Name()))
			src = b.newVar(newSrcName, src.Type())
			stmts = append(stmts, assignStmt)
		}
	}

	if b.buildConfig.BuildMode == BuildModeConv {
		if ret, ok := b._shallowCopy(dst, src); ok {
			return append(stmts, ret...)
		}
	}
	// src pointer to
	if srcPtr, ok := src.Type().(*types.Pointer); ok {
		return append(stmts, b.dePointerSrcStmt(dst, src, srcPtr.Elem())...)
	}
	switch dstType := dst.Type().(type) {
	case *types.Pointer:
		_, srcIsStruct := isStruct(src.Type())
		// check has generated func
		if !b.rootNode && srcIsStruct && isPointerToStruct(dstType) {
			funcName := b.GenFuncName(types.NewPointer(src.Type()), dstType, b.buildConfig)
			convSrcName := func() string {
				return addressName(src.Name(), 1)
			}()
			assignStmt := buildAssignStmt(dst.Name(), fmt.Sprintf("%s(%s)", funcName, convSrcName))
			if _, ok := b.genFunc[funcName]; !ok {
				curGenFunc := b.curGenFunc
				b.BuildFunc(dstType, types.NewPointer(src.Type()), b.buildConfig)
				b.curGenFunc = curGenFunc
			}
			stmts = append(stmts, assignStmt)
			return stmts
		}
		dstName := ptrToName(dst.Name(), 1)
		dstElemVar := b.newVar(dstName, dstType.Elem())
		elementStmt := b.buildStmt(dstElemVar, src)
		// dst = new(dst.Type)
		initAssign := buildAssignStmt(dst.Name(), fmt.Sprintf("new(%s)", b.importer.ImportType(dstElemVar.Type())))
		stmts = append(stmts, initAssign)
		stmts = append(stmts, elementStmt...)
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
			if ss, ok := b._shallowCopy(dst, src); ok {
				return append(stmts, ss...)
			}
		}
		dstUnderVar := b.newVar(dst.Name(), dstUnderType)
		b.varToNamed[dstUnderVar] = dstType
		return b.buildStmt(dstUnderVar, src)
	case *types.Struct:
		srcStructType, ok := isStruct(src.Type())
		if !ok {
			b.logger.Notice("omit %s :%s type is not a struct/pointer to struct", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}
		srcName := strings.TrimPrefix(src.Name(), "*")
		dstName := strings.TrimPrefix(dst.Name(), "*") // for struct type, compiler can de Pointer automatically
		var dstStructName string
		if fromNamed, has := b.varToNamed[dst]; has {
			dstStructName = fromNamed.String()
		} else {
			dstStructName = dstType.String()
		}
	OUTER:
		for i := range dstType.NumFields() {
			dstField := dstType.Field(i)
			if !dstField.Exported() {
				continue
			}
			dstFieldName := dstField.Name()
			if dstField.Embedded() {
				dstVar := b.newVar(dstName+"."+dstFieldName, dstField.Type())
				fieldStmt := b.buildStmt(dstVar, src)
				stmts = append(stmts, fieldStmt...)
				continue
			}
			// check ignore
			b.fieldPath.Push(fieldStep{
				dst: field{
					name:       dstFieldName,
					structName: dstStructName,
				},
			})
			for _, ignoreType := range b.buildConfig.Ignore {
				if ignoreType.IgnoreSide == SideSrc {
					continue
				}
				if b.fieldPath.matchIgnore(ignoreType, nil, dstField.Type()) {
					b.logger.Printf("apply ignore on dst %s", src.Name())
					b.buildCommentExpr(&stmts, "apply ignore option on %s", dstName+"."+dstFieldName)
					continue OUTER
				}
			}
			// match srcField
			b.fieldPath.Pop()
			if srcField, ok := b.matchField(dstField, srcStructType, src.Type().String()); ok {
				dstVarName := dstName + "." + dstFieldName
				srcVarName := srcName + "." + srcField.Name()
				b.logger.Printf("assign [%s(%s) -> %s(%s)]", srcVarName, srcField.Type().String(),
					dstVarName, dstField.Type().String())
				b.fieldPath.Push(fieldStep{
					src: field{
						name:       srcField.Name(),
						structName: src.Type().String(),
					},
					dst: field{
						name:       dstField.Name(),
						structName: dstStructName,
					},
				})
				dstVar := b.newVar(dstVarName, dstField.Type())
				srcVar := b.newVar(srcVarName, srcField.Type())
				fieldStmt := b.buildStmt(dstVar, srcVar)
				stmts = append(stmts, fieldStmt...)
				b.fieldPath.Pop()
			} else {
				b.logger.Notice("omit %s :not find match field in %s", dstFieldName, srcName)
				b.buildCommentExpr(&stmts, "omit "+dstFieldName)
			}
		}
		return stmts
	case *types.Array:
		srcArrType, isSlice, ok := convSliceToArray(src.Type())
		if !ok {
			b.logger.Notice("omit %s :%s type is not a array/slice", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}

		var srcName = src.Name()
		// for i := 0; i<n ; i++ {}
		b.pushScope("range@" + srcName)
		symbol := b.scope.NextSymbol("i")
		rangeIndex := symbol.Name
		forStmt := &ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent(rangeIndex)},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{ast.NewIdent("0")},
			},
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent(rangeIndex),
				Op: token.LSS,
				Y:  ast.NewIdent(strconv.FormatInt(dstType.Len(), 10)),
			},
			Post: &ast.IncDecStmt{
				X:   ast.NewIdent(rangeIndex),
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{},
		}
		if isSlice || srcArrType.Len() > dstType.Len() {
			forStmt.Cond = &ast.BinaryExpr{
				X:  forStmt.Cond,
				Op: token.LAND,
				Y: &ast.BinaryExpr{
					X:  ast.NewIdent(rangeIndex),
					Op: token.LSS,
					Y: &ast.CallExpr{
						Fun:  ast.NewIdent("len"),
						Args: []ast.Expr{ast.NewIdent(srcName)},
					},
				},
			}
		}
		dstElemVar := b.newVar(parenthesesName(
			dst.Name())+fmt.Sprintf("[%s]", rangeIndex), dstType.Elem())
		srcElemVar := b.newVar(
			srcName+fmt.Sprintf("[%s]", rangeIndex), srcArrType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		b.popScope()
		forStmt.Body.List = append(forStmt.Body.List, elementStmt...)
		stmts = append(stmts, forStmt)
		return stmts
	case *types.Map:
		srcType, ok := src.Type().Underlying().(*types.Map)
		if !ok {
			b.logger.Notice("omit %s :%s type is not a map", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
		}
		var srcName = src.Name()
		ifStmt := buildIfStmt(fmt.Sprintf("len(%s)", srcName), token.GTR, "0")

		dstKeyTypeStr := b.importer.ImportType(dstType.Key())
		dstValueTypeStr := b.importer.ImportType(dstType.Elem())

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
		b.pushScope("range@" + srcName)
		tmpK, tmpV := b.scope.NextPair("tmpK", "tmpV")
		k, v := b.scope.NextPair("k", "v")
		dstKeyVar := b.newVar(tmpK.Name, dstType.Key())
		dstValueVar := b.newVar(tmpV.Name, dstType.Elem())
		srcKeyVar := b.newVar(k.Name, srcType.Key())
		srcValueVar := b.newVar(v.Name, srcType.Elem())
		rangeStmt := &ast.RangeStmt{
			Key:   ast.NewIdent(k.Name),
			Value: ast.NewIdent(v.Name),
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
		b.popScope()
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignKStmt...)
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignVStmt...)
		assignMapStmt := buildAssignStmt(fmt.Sprintf("%s[%s]", parenthesesName(dst.Name()), dstKeyVar.Name()), dstValueVar.Name())
		rangeStmt.Body.List = append(rangeStmt.Body.List, assignMapStmt)
		stmts = append(stmts, ifStmt)
		return stmts
	case *types.Slice:
		srcSliceType, isArray, ok := convArrayToSlice(src.Type())
		if !ok {
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
			b.logger.Notice("omit %s :%s type is not a slice/array", dst.Name(), src.Name())
			b.buildCommentExpr(&stmts, "omit "+dst.Name())
			return stmts
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
		b.pushScope("range@" + srcName)
		symbol := b.scope.NextSymbol("i")
		rangeIndex := symbol.Name
		forStmt := &ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent(rangeIndex)},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{ast.NewIdent("0")},
			},
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent(rangeIndex),
				Op: token.LSS,
				Y: &ast.CallExpr{
					Fun:  ast.NewIdent("len"),
					Args: []ast.Expr{ast.NewIdent(srcName)},
				},
			},
			Post: &ast.IncDecStmt{
				X:   ast.NewIdent(rangeIndex),
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{},
		}
		dstElemVar := b.newVar(parenthesesName(
			dst.Name())+fmt.Sprintf("[%s]", rangeIndex), dstType.Elem())
		srcElemVar := b.newVar(
			srcName+fmt.Sprintf("[%s]", rangeIndex), srcSliceType.Elem())
		elementStmt := b.buildStmt(dstElemVar, srcElemVar)
		b.popScope()
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
		if ss, ok := b._shallowCopy(dst, src); ok {
			return append(stmts, ss...)
		}
		b.logger.Notice("omit %s :basic type can't cast from %s (or it pointers to)", dst.Name(), src.Name())
		b.buildCommentExpr(&stmts, "omit "+dst.Name())
	default:
		b.logger.Notice("omit %s :type not support yet", dst.Name())
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
	case *types.Struct, *types.Basic, *types.Pointer:
		needParentheses = false
	}
	ptrToSrcName := ptrToName(src.Name(), 1)
	if needParentheses {
		ptrToSrcName = parenthesesName(ptrToSrcName)
	}
	srcElemVar := b.newVar(ptrToSrcName, srcElemType)
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

func isStruct(v types.Type) (structType *types.Struct, ok bool) {
	if ret, ok := v.(*types.Struct); ok {
		return ret, true
	}
	// check if src is a Named struct
	if namedTypes, ok := v.(*types.Named); ok {
		if ret, ok := namedTypes.Underlying().(*types.Struct); ok {
			return ret, true
		}
	}
	return nil, false
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
		b.f.Decls = append(b.f.Decls, b.genFunc[name].GenFunc)
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
	sb.WriteString("// Code generated by github.com/ycl2018/go-conv DO NOT EDIT.\n// +build !goconv_gen\n\n")
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
	var funcName string
	srcTypeName, dstTypeName := b.importer.ImportType(src), b.importer.ImportType(dst)
	switch buildConfig.BuildMode {
	case BuildModeCopy:
		funcName = "Copy" + cleanName(srcTypeName) + "To" + cleanName(dstTypeName)
	default:
		funcName = cleanName(srcTypeName) + "To" + cleanName(dstTypeName)
	}
	for i := 0; ; i++ {
		if i > 0 {
			funcName = funcName + strconv.Itoa(i)
		}
		if gf, ok := b.genFunc[funcName]; !ok {
			return funcName
		} else if gf.buildConfig.String() == buildConfig.String() {
			return funcName
		}
	}
}

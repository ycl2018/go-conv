package internal

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func ParseVarsToConv(pkgs []*packages.Package) (map[*Package][]*ConvVar, error) {

	var varsToConv = map[*Package][]*ConvVar{}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("[go-conv] parse err in %s \n", pkg.PkgPath))
			for _, e := range pkg.Errors {
				sb.WriteString(fmt.Sprintf("%s\t%s\n", e.Pos, e.Msg))
			}
			return nil, errors.New(sb.String())
		}
		var p = &Package{Package: pkg}
		cp := &CommentParser{pkg: pkg}
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				// only resolve the package level Decl
				if gd, ok := decl.(*ast.GenDecl); ok {
					var isTarget = false
					if gd.Doc == nil {
						continue
					}
					for _, comment := range gd.Doc.List {
						if strings.Contains(comment.Text, "go-conv:generate") {
							isTarget = true
							break
						}
					}
					if !isTarget {
						continue
					}
					buildConfig := DefaultBuildConfig()
					err := cp.Parse(astFile, gd.Doc, &buildConfig)
					if err != nil {
						return nil, fmt.Errorf("[go-conv] parse comment err:%w", err)
					}
					for _, spec := range gd.Specs {
						if vs, ok := spec.(*ast.ValueSpec); ok {
							if vs.Doc != nil {
								buildConfig = buildConfig.Clone()
								err = cp.Parse(astFile, vs.Doc, &buildConfig)
								if err != nil {
									return nil, fmt.Errorf("[go-conv] parse comment err:%w", err)
								}
							}
							tye := pkg.TypesInfo.TypeOf(vs.Type)
							if sig, ok := tye.(*types.Signature); ok {
								if sig.Params().Len() == 0 || sig.Results().Len() == 0 {
									return nil, fmt.Errorf(
										"[go-conv] err: 0 params/results func Signature found at\n%s",
										pkg.Fset.Position(vs.Pos()).String())
								}
								// get package dir
								if p.Dir == "" {
									f := pkg.Fset.File(spec.Pos())
									fileName := f.Name()
									p.Dir = filepath.Dir(fileName)
								}
								varsToConv[p] = append(varsToConv[p], &ConvVar{
									VarSpec:     vs,
									Signature:   sig,
									BuildConfig: buildConfig,
								})
							}
						}
					}
				}
			}
		}
	}
	return varsToConv, nil
}

const (
	applyOptionsType           = "[]github.com/ycl2018/go-conv/option.Option"
	ignoreFieldsOption         = "WithIgnoreFields"
	ignoreTypesOption          = "WithIgnoreTypes"
	transformerOption          = "WithTransformer"
	filterOption               = "WithFilter"
	noInitOption               = "WithNoInitFunc"
	fieldMatchOption           = "WithFieldMatch"
	fieldCaseInsensitiveOption = "WithMatchCaseInsensitive"
)

type CommentParser struct {
	pkg *packages.Package
}

func (c CommentParser) Parse(astFile *ast.File, doc *ast.CommentGroup, config *BuildConfig) error {
	if doc == nil {
		return nil
	}
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "go-conv:copy") {
			config.BuildMode = BuildModeCopy
		} else if strings.Contains(comment.Text, "go-conv:conv") {
			config.BuildMode = BuildModeConv
		} else if strings.Contains(comment.Text, "go-conv:apply") {
			err := c.parseApply(astFile, comment, config)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c CommentParser) parseApply(astFile *ast.File, comment *ast.Comment, ret *BuildConfig) error {
	fields := strings.Fields(strings.TrimPrefix(comment.Text, "//"))
	if len(fields) <= 1 {
		return fmt.Errorf("%s:not set apply value", c.pkg.Fset.Position(comment.Slash))
	}
	applyValueName := fields[1]
	applyValue := c.pkg.Types.Scope().Lookup(applyValueName)
	if applyValue == nil {
		return fmt.Errorf("%s:not find apply value of name:%s",
			c.pkg.Fset.Position(comment.Slash), applyValueName)
	}
	typeString := applyValue.Type().String()
	if typeString != applyOptionsType {
		return fmt.Errorf("%s:not Option Slice", c.pkg.Fset.Position(comment.Slash))
	}
	nodes, _ := astutil.PathEnclosingInterval(astFile, applyValue.Pos(), applyValue.Pos())

	for _, node := range nodes {
		vs, ok := node.(*ast.ValueSpec)
		if !ok || len(vs.Values) == 0 {
			continue
		}
		compositeLit, ok := vs.Values[0].(*ast.CompositeLit)
		if !ok {
			continue
		}

		for _, elt := range compositeLit.Elts {
			callExpr, ok := elt.(*ast.CallExpr)
			if !ok {
				continue
			}
			optionFn := callExpr.Fun.(*ast.SelectorExpr).Sel.Name
			switch optionFn {
			case ignoreFieldsOption:
				var ignoreFields []string
				ast.Inspect(callExpr.Args[1], func(n ast.Node) bool {
					if basicLit, ok := n.(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
						ignoreFields = append(ignoreFields, strings.Trim(basicLit.Value, "\""))
					}
					return true
				})
				var paths []string
				if len(callExpr.Args) > 2 {
					for i := 2; i < len(callExpr.Args); i++ {
						paths = append(paths, strings.Trim(callExpr.Args[i].(*ast.BasicLit).Value, "\""))
					}
				}
				structType := c.pkg.TypesInfo.TypeOf(callExpr.Args[0])
				elemType, _, _ := dePointer(structType)
				ret.Ignore = append(ret.Ignore, IgnoreType{
					Tye:    elemType.String(),
					Fields: ignoreFields,
					Paths:  paths,
				})
				DefaultLogger.Printf("[go-conv] find comment on %s: config ignore %s fields: %v with paths:%v",
					c.pkg.Fset.Position(elt.Pos()), elemType, ignoreFields, paths)
			case ignoreTypesOption:
				var paths []string
				if len(callExpr.Args) > 1 {
					for i := 1; i < len(callExpr.Args); i++ {
						paths = append(paths, strings.Trim(callExpr.Args[i].(*ast.BasicLit).Value, "\""))
					}
				}
				ignoreType := c.pkg.TypesInfo.TypeOf(callExpr.Args[0])
				elemType, _, _ := dePointer(ignoreType)
				ret.Ignore = append(ret.Ignore, IgnoreType{
					Tye:   elemType.String(),
					Paths: paths,
				})
				DefaultLogger.Printf("[go-conv] find comment on %s: config ignore type:%s with paths:%v",
					c.pkg.Fset.Position(elt.Pos()), elemType, paths)
			case transformerOption:
				transferFuncName, ok := callExpr.Args[0].(*ast.Ident)
				if !ok {
					return fmt.Errorf("%s:should be a funcName, not support anonymous function",
						c.pkg.Fset.Position(callExpr.Args[0].Pos()))
				}
				transfer, ok := c.pkg.TypesInfo.TypeOf(callExpr.Args[0]).(*types.Signature)
				if !ok {
					return fmt.Errorf("%s:%s shoule be signature func(T)V", transferFuncName,
						c.pkg.Fset.Position(callExpr.Pos()))
				}
				if transfer.Params().Len() != 1 || transfer.Results().Len() != 1 {
					return fmt.Errorf("%s:%s shoule be signature func(T)V", transferFuncName,
						c.pkg.Fset.Position(callExpr.Pos()))
				}
				from, to := transfer.Params().At(0).Type(), transfer.Results().At(0).Type()
				var paths []string
				if len(callExpr.Args) > 1 {
					for i := 1; i < len(callExpr.Args); i++ {
						paths = append(paths, strings.Trim(callExpr.Args[i].(*ast.BasicLit).Value, "\""))
					}
				}
				ret.Transfer = append(ret.Transfer, Transfer{
					From:     from.String(),
					To:       to.String(),
					FuncName: transferFuncName.Name,
					Paths:    paths,
				})
				DefaultLogger.Printf("[go-conv] find comment on %s: config transfer %s from:%s to %s with "+
					"paths:%s",
					c.pkg.Fset.Position(elt.Pos()), transferFuncName, from.String(), to.String(), paths)
			case filterOption:
				filterFuncName, ok := callExpr.Args[0].(*ast.Ident)
				if !ok {
					return fmt.Errorf("%s:should be a funcName, not support anonymous function",
						c.pkg.Fset.Position(callExpr.Args[0].Pos()))
				}
				transfer, ok := c.pkg.TypesInfo.TypeOf(callExpr.Args[0]).(*types.Signature)
				if !ok {
					return fmt.Errorf("%s:%s shoule be signature func(T)T", filterFuncName,
						c.pkg.Fset.Position(callExpr.Pos()))
				}
				if transfer.Params().Len() != 1 || transfer.Results().Len() != 1 {
					return fmt.Errorf("%s:%s shoule be signature func(T)T", filterFuncName,
						c.pkg.Fset.Position(callExpr.Pos()))
				}
				from, to := transfer.Params().At(0).Type(), transfer.Results().At(0).Type()
				if from.String() != to.String() {
					return fmt.Errorf("%s:%s shoule be signature func(T)T", filterFuncName,
						c.pkg.Fset.Position(callExpr.Pos()))
				}
				var paths []string
				if len(callExpr.Args) > 1 {
					for i := 1; i < len(callExpr.Args); i++ {
						paths = append(paths, strings.Trim(callExpr.Args[i].(*ast.BasicLit).Value, "\""))
					}
				}
				ret.Filter = append(ret.Filter, Filter{
					Typ:      from.String(),
					FuncName: filterFuncName.Name,
					Paths:    paths,
				})
				DefaultLogger.Printf("[go-conv] find comment on %s: config filter %s on %s with paths:%s",
					c.pkg.Fset.Position(elt.Pos()), filterFuncName, from.String(), paths)
			case noInitOption:
				ret.NoInit = true
				DefaultLogger.Printf("[go-conv] find comment on %s: config no generate init func",
					c.pkg.Fset.Position(elt.Pos()))
			case fieldMatchOption:
				structType := c.pkg.TypesInfo.TypeOf(callExpr.Args[0])
				elemType, _, _ := dePointer(structType)
				elemTypeStr := elemType.String()
				cl, ok := callExpr.Args[1].(*ast.CompositeLit)
				if !ok {
					return fmt.Errorf("%s:type shoule be map[string]string literal",
						c.pkg.Fset.Position(callExpr.Args[1].Pos()))
				}
				for _, expr := range cl.Elts {
					kvExpr := expr.(*ast.KeyValueExpr)
					k, ok1 := kvExpr.Key.(*ast.BasicLit)
					v, ok2 := kvExpr.Value.(*ast.BasicLit)
					if !ok1 || !ok2 {
						return fmt.Errorf("%s:type shoule be map[string]string literal",
							c.pkg.Fset.Position(kvExpr.Key.Pos()))
					}
					from, to := strings.Trim(k.Value, "\""), strings.Trim(v.Value, "\"")
					ret.FieldMatcher.AddMatch(elemTypeStr, from, to)
					DefaultLogger.Printf("[go-conv] find comment on %s: config type %s match field from %s "+
						"to %s",
						c.pkg.Fset.Position(elt.Pos()), elemTypeStr, from, to)
				}
			case fieldCaseInsensitiveOption:
				ret.CaseInsensitive = true
				DefaultLogger.Printf("[go-conv] find comment on %s: config field matched by case-insensitive",
					c.pkg.Fset.Position(elt.Pos()))
			default:
				panic("expect unreachable")
			}
		}

	}
	return nil
}

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
					buildConfig, err := cp.Parse(astFile, gd.Doc)
					if err != nil {
						return nil, fmt.Errorf("[go-conv] parse comment err:%w", err)
					}
					for _, spec := range gd.Specs {
						if vs, ok := spec.(*ast.ValueSpec); ok {
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
	applyOptionsType  = "[]github.com/ycl2018/go-conv/option.Option"
	ignoreFieldsFunc  = "WithIgnoreFields"
	ignoreTypesFunc   = "WithIgnoreTypes"
	ignoreIndexesFunc = "WithIgnoreIndexes"
	ignoreKeysFunc    = "WithIgnoreKeys"
)

type CommentParser struct {
	pkg *packages.Package
}

func (c CommentParser) Parse(astFile *ast.File, doc *ast.CommentGroup) (BuildConfig, error) {
	var ret = DefaultBuildConfig
	if doc == nil {
		return ret, nil
	}
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "go-conv:copy") {
			ret.BuildMode = BuildModeCopy
		} else if strings.Contains(comment.Text, "go-conv:conv") {
			ret.BuildMode = BuildModeConv
		} else if strings.Contains(comment.Text, "go-conv:apply") {
			err := c.parseApply(astFile, comment, &ret)
			if err != nil {
				return ret, err
			}
		}
	}
	return ret, nil
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
	DefaultLogger.Notice("%s", nodes)

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
			DefaultLogger.Printf("%s", optionFn)
			switch optionFn {
			case ignoreFieldsFunc:
				structType := c.pkg.TypesInfo.TypeOf(callExpr.Args[0])
				var ignoreFields []string
				ast.Inspect(callExpr.Args[1], func(n ast.Node) bool {
					if basicLit, ok := n.(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
						ignoreFields = append(ignoreFields, basicLit.Value)
					}
					return true
				})
				DefaultLogger.Printf("[go-conv] find comment on %s: config ignore %s fields: %v",
					c.pkg.Fset.Position(elt.Pos()), structType, ignoreFields)
				ret.Ignore[IgnoreType{
					typ:  structType,
					kind: IgnoreStructFields,
				}] = ignoreFields
			case ignoreTypesFunc:
				for _, arg := range callExpr.Args {
					ignoreType := c.pkg.TypesInfo.TypeOf(arg)
					ret.Ignore[IgnoreType{
						typ:  ignoreType,
						kind: IgnoreTypes,
					}] = struct{}{}
					DefaultLogger.Printf("[go-conv] find comment on %s: config ignore type:%s",
						c.pkg.Fset.Position(elt.Pos()), ignoreType)
				}
			}
		}

	}
	return nil
}

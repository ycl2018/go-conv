package internal

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func ParseVarsToConv(pkgs []*packages.Package) (map[*Package][]*ConvVar, error) {

	var varsToConv = map[*Package][]*ConvVar{}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("[go-conv] parse %s err:\n%v", pkg.PkgPath, pkg.Errors)
		}
		var p = &Package{Package: pkg}
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
					var buildConfig = parseConfigFromComment(gd.Doc)
					for _, spec := range gd.Specs {
						if vs, ok := spec.(*ast.ValueSpec); ok {
							tye := pkg.TypesInfo.TypeOf(vs.Type)
							if sig, ok := tye.(*types.Signature); ok {
								if sig.Params().Len() == 0 || sig.Results().Len() == 0 {
									return nil, fmt.Errorf(
										"[go-conv] err: 0 params/results func Signature found at %s",
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

func parseConfigFromComment(doc *ast.CommentGroup) BuildConfig {
	var ret = DefaultBuildConfig
	if doc == nil {
		return ret
	}
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "go-conv:copy") {
			ret.BuildMode = BuildModeCopy
		} else if strings.Contains(comment.Text, "go-conv:conv") {
			ret.BuildMode = BuildModeConv
		}
	}
	return ret
}

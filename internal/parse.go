package internal

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"path/filepath"
	"strings"
)

func ParseVarsToConv(pkgs []*packages.Package) (map[*Package][]*ConvVar, error) {

	var varsToConv = map[*Package][]*ConvVar{}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("[go-conv] package:%s contain syntax errors: %v", pkg.PkgPath, pkg.Errors)
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
									VarSpec:   vs,
									Signature: sig,
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

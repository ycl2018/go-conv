package internal

import (
	"go/types"
	"strconv"
)

type Importer struct {
	curPkgPath      string
	pkgToName       map[string]string
	importedPkgName map[string]int
	imported        []*types.Package
}

func NewImporter(curPkgPath string) *Importer {
	return &Importer{
		curPkgPath:      curPkgPath,
		pkgToName:       map[string]string{},
		importedPkgName: map[string]int{},
	}
}

func (i *Importer) ImportType(t types.Type) string {
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
	resolve(t)
	if pkgPath == "" || pkgPath == i.curPkgPath {
		return typPrefix + typeName
	}
	pkgImportName, ok := i.pkgToName[pkgPath]
	if !ok {
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

	name := typPrefix + pkgImportName + "." + typeName
	return name
}

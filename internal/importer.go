package internal

import (
	"fmt"
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

	var typeName string
	var typPrefix string
	var pkg *types.Package
	var resolve func(tye types.Type)
	resolve = func(tye types.Type) {
		switch varType := tye.(type) {
		case *types.Named:
			pkg = varType.Obj().Pkg()
			typeName += varType.Obj().Name()
			return
		case *types.Basic:
			typeName += varType.Name()
			if varType.Kind() == types.UnsafePointer {
				pkg = types.NewPackage("unsafe", "unsafe")
			}
			return
		case *types.Slice:
			typPrefix += "[]"
			resolve(varType.Elem())
		case *types.Array:
			typPrefix += fmt.Sprintf("[%d]", varType.Len())
			resolve(varType.Elem())
		case *types.Map:
			keyName, vName := i.ImportType(varType.Key()), i.ImportType(varType.Elem())
			typPrefix += fmt.Sprintf("map[%s]%s", keyName, vName)
		case *types.Pointer:
			typPrefix += "*"
			resolve(varType.Elem())
		default:
			panic("expect unreachable")
		}
	}
	resolve(t)
	var pkgPath string
	var pkgName string
	if pkg != nil {
		pkgPath = pkg.Path()
		pkgName = pkg.Name()
	}
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

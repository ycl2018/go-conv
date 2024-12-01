package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"path/filepath"

	. "go-conv/internal"
)

// flags
var (
	dryRun  = flag.Bool("n", false, "dry run: show changes, but don't apply them")
	verbose = flag.Bool("v", false, "verbose: enable verbose output log")
)

// seams for testing
var (
	stderr    io.Writer = os.Stderr
	writeFile           = os.WriteFile
)

const usage = `go-conv: generate golang types convert/copy source code.

Usage: go-conv [flags] package...

The package... arguments specify a list of packages
in the style of the go tool; see "go help packages".
Hint: use "all" or "..." to match the entire workspace.

Flags:
  -n:	       dry run: show generate code, but don't write it to file
  -v:		   verbose: enable verbose output log
`

func main() {
	flag.Parse()
	flag.Usage()
	if len(flag.Args()) == 0 {
		fmt.Fprint(stderr, usage)
		os.Exit(1)
	}

	if err := generate(flag.Args()...); err != nil {
		fmt.Fprintf(stderr, err.Error())
		os.Exit(1)
	}
}

const parserLoadMode = packages.NeedName | packages.NeedImports | packages.NeedDeps |
	packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo

func generate(patterns ...string) error {

	loadConf := &packages.Config{
		Mode: parserLoadMode,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			return parser.ParseFile(fset, filename, src, parser.SkipObjectResolution|parser.ParseComments)
		},
	}

	pkgs, err := packages.Load(loadConf, patterns...)
	if err != nil {
		return fmt.Errorf("[go-conv] loading packages err:\n%w", err)
	}

	varsToConv, err := ParseVarsToConv(pkgs)
	if err != nil {
		return err
	}

	if len(varsToConv) == 0 {
		return fmt.Errorf("[go-conv] not find valid function to convert in path:%s", patterns)
	}

	for p, vars := range varsToConv {
		fileToGen := &ast.File{
			Name: ast.NewIdent(p.Name),
		}
		impt := NewImporter()
		builder := NewBuilder(fileToGen, p.Types, impt)
		for _, v := range vars {
			src, dst := v.Signature.Params().At(0), v.Signature.Results().At(0)
			builder.BuildFunc(dst.Type(), src.Type())
		}
		content, err := builder.Generate()
		if err != nil {
			return err
		}
		genFileName := "goconv_gen.go"
		// write file
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("[go-conv] get working dir err:\n%w", err)
		}
		fileName := filepath.Join(filepath.Dir(dir), p.PkgPath, genFileName)
		if *dryRun {
			fmt.Fprintf(stderr, "[go-conv] generate to %s content:\n\n%s", fileName, content)
		} else {
			err = os.WriteFile(fileName, content, 0644)
			if err != nil {
				return fmt.Errorf("[go-conv] write file %s err:%w", fileName, err)
			}
		}
	}
	return nil
}

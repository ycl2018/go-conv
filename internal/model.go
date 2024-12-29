package internal

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// ConvVar vars to generate from
type ConvVar struct {
	VarSpec     *ast.ValueSpec
	Signature   *types.Signature
	BuildConfig BuildConfig
}

type Package struct {
	*packages.Package
	Dir string
}

type Transfer struct {
	From, To string
	FuncName string
	Paths    []string
}

type Filter struct {
	Typ      string
	FuncName string
	Paths    []string
}

type BuildConfig struct {
	BuildMode BuildMode
	Ignore    []IgnoreType
	Transfer  []Transfer
	Filter    []Filter
}

type IgnoreType struct {
	Tye    string
	Fields []string
	Paths  []string
}

func (b BuildConfig) String() string {
	return fmt.Sprintf("BuildConfig<BuildMode: %s>", b.BuildMode)
}

var DefaultBuildConfig = BuildConfig{
	BuildMode: BuildModeConv,
}

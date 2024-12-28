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
	typ      string
	FuncName string
	Paths    []string
}

type BuildConfig struct {
	BuildMode BuildMode
	Ignore    []IgnoreType
	Transfer  []Transfer
	Filter    []Filter
}

type IgnoreKind int

const (
	IgnoreStructFields = iota + 1
	IgnoreSliceIndexes
	IgnoreMapKeys
	IgnoreTypes
)

type IgnoreType struct {
	typ    string
	fields []string
}

func (b BuildConfig) String() string {
	return fmt.Sprintf("BuildConfig<BuildMode: %s>", b.BuildMode)
}

var DefaultBuildConfig = BuildConfig{
	BuildMode: BuildModeConv,
}

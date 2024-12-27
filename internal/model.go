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

type BuildConfig struct {
	BuildMode BuildMode
	Ignore    map[IgnoreType]any
}

type IgnoreKind int

const (
	IgnoreStructFields = iota + 1
	IgnoreSliceIndexes
	IgnoreMapKeys
	IgnoreTypes
)

type IgnoreType struct {
	typ  types.Type
	kind IgnoreKind
}

func (b BuildConfig) String() string {
	return fmt.Sprintf("BuildConfig<BuildMode: %s>", b.BuildMode)
}

var DefaultBuildConfig = BuildConfig{
	BuildMode: BuildModeConv,
	Ignore:    map[IgnoreType]any{},
}

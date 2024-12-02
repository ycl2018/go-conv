package internal

import (
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
}

var DefaultBuildConfig = BuildConfig{
	BuildMode: BuildModeConv,
}

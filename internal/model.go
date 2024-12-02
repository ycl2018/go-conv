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
}

func (b BuildConfig) String() string {
	return fmt.Sprintf("BuildConfig<BuildMode: %s>", b.BuildMode)
}

var DefaultBuildConfig = BuildConfig{
	BuildMode: BuildModeConv,
}

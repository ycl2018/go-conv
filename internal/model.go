package internal

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
)

// ConvVar vars to generate from
type ConvVar struct {
	VarSpec   *ast.ValueSpec
	Signature *types.Signature
	Package   *packages.Package
}

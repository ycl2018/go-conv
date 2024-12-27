package cases

import (
	"github.com/ycl2018/go-conv/option"
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// BasicToBasic conv a Basic to b Basic
// go-conv:generate
// go-conv:apply basicConvOpts
var (
	BasicToBasic func(*a.Basic) *b.Basic
)

var basicConvOpts = []option.Option{
	option.WithIgnoreFields(a.Basic{}, []string{"Int64", "Int32"}),
}

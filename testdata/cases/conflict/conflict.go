package conflict

import (
	"github.com/ycl2018/go-conv/option"
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var BasicToBasic func(*a.Basic) *b.Basic

// go-conv:generate
// go-conv:apply basicConvOpts
var BasicToBasicOmit2 func(*a.Basic) *b.Basic

var basicConvOpts = []option.Option{
	option.WithIgnoreFields(a.Basic{}, []string{"Uint8", "Uint16"}),
}

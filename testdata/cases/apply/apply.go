package cases

import (
	"github.com/ycl2018/go-conv/option"
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
	"strconv"
)

// BasicToBasic conv a Basic to b Basic
// go-conv:generate
// go-conv:apply basicConvOpts
var (
	BasicToBasic func(p *a.Struct) *b.Struct
)

var basicConvOpts = []option.Option{
	option.WithIgnoreFields(a.Struct{}, []string{"Student"}),
	option.WithTransformer(transfer, "Student2.Class.Grade"),
	option.WithFilter(filter, "Student2.Teachers"),
}

func transfer(t int) string {
	return strconv.Itoa(t)
}

func filter(arr []string) []string {
	return arr
}

package cases

import (
	"strconv"

	"github.com/ycl2018/go-conv/option"
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// Struct2Struct conv a Basic to b Basic
// go-conv:generate
// go-conv:apply basicConvOpts
var (
	Struct2Struct func(p *a.Struct) *b.Struct
)

var basicConvOpts = []option.Option{
	option.WithIgnoreFields(a.Struct{}, []string{"Pojo"}),
	option.WithIgnoreDstFields(b.Struct{}, []string{"IgnoreField"}),
	option.WithIgnoreTypes(a.Student{}, "Student3"),
	option.WithIgnoreDstTypes(b.Pojo{}, "IgnoreType"),
	option.WithTransformer(transfer, "Student2.Class.Grade"),
	option.WithFilter(filter, "Student2.Teachers"),
	option.WithFieldMatch(a.Struct{}, map[string]string{
		"Match": "Match_",
	}),
	option.WithMatchCaseInsensitive(),
	option.WithNoInitFunc(),
}

func transfer(t int) string {
	return strconv.Itoa(t)
}

func filter(arr []string) []string {
	return arr
}

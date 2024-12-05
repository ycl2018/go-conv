package generic

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var (
	GenToGen      func(generic *a.Generic[string, string]) *b.Generic[string, string]
	GenToGenCast  func(generic *a.Generic[string, int]) *b.Generic[string, int64]
	GenToGenSlice func(generic *a.Generic[string, []int]) *b.Generic[string, []int64]
)

package cast

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var (
	IntToUint     func(cast *a.Cast[int]) *b.Cast[uint]
	IntToFloat    func(cast *a.Cast[int]) *b.Cast[float32]
	FloatToInt    func(cast a.Cast[float32]) *b.Cast[int]
	FloatToIntPtr func(cast a.Cast[float32]) *b.Cast[*int]
)

package cast

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var (
	IntToUint         func(cast *a.Cast[int]) *b.Cast[uint]
	IntToFloat        func(cast *a.Cast[int]) *b.Cast[float32]
	FloatToInt        func(cast a.Cast[float32]) *b.Cast[int]
	FloatToIntPtr     func(cast a.Cast[float32]) *b.Cast[*int]
	StringToByteSlice func(cast *a.Cast[string]) *b.Cast[[]byte]
	ByteSliceToString func(cast a.Cast[[]byte]) *b.Cast[string]
	StringToRuneSlice func(cast *a.Cast[string]) *b.Cast[[]rune]
	RuneSliceToString func(cast a.Cast[[]rune]) *b.Cast[[]rune]
)

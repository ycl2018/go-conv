package nested

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var (
	StructPtrSliceToStructPtrSlice func(src []*a.Foo) (dst []*b.Foo)
	NestedSliceToNestedSlice       func(src *a.NestedSlice) (dst *b.NestedSlice)
	NestedSliceToNestedSlice2      func(src *b.NestedSlice) (dst *a.NestedSlice)
)

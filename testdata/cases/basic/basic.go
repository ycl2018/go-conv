package basic

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var (
	BasicToBasic               func(*a.Basic) *b.Basic
	BasicToBasicPtr            func(a.Basic) *b.Basic
	BasicPtrToBasic            func(*a.Basic) b.Basic
	BasicToBasicPointer        func(*a.Basic) *b.BasicPtr
	BasicPointerToBasic        func(ptr *a.BasicPtr) *b.Basic
	BasicPointerToBasicPointer func(ptr *a.BasicPtr) *b.BasicPtr
)

// go-conv:generate
// go-conv:copy
var (
	CopyBasicToBasic               func(*a.Basic) *b.Basic
	CopyBasicToBasicPointer        func(*a.Basic) *b.BasicPtr
	CopyBasicPointerToBasic        func(ptr *a.BasicPtr) *b.Basic
	CopyBasicPointerToBasicPointer func(ptr *a.BasicPtr) *b.BasicPtr
)

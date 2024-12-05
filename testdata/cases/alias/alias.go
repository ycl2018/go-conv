package alias

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var (
	BasicToBasic               func(*a.BasicAlias) *b.BasicAlias
	BasicToBasicPointer        func(*a.BasicAlias) *b.BasicAliasPtr
	BasicPointerToBasic        func(ptr *a.BasicAliasPtr) *b.BasicAlias
	BasicPointerToBasicPointer func(ptr *a.BasicAliasPtr) *b.BasicAliasPtr
)

// go-conv:generate
// go-conv:copy
var (
	CopyBasicToBasic               func(*a.BasicAlias) *b.BasicAlias
	CopyBasicToBasicPointer        func(*a.BasicAlias) *b.BasicAliasPtr
	CopyBasicPointerToBasic        func(ptr *a.BasicAliasPtr) *b.BasicAlias
	CopyBasicPointerToBasicPointer func(ptr *a.BasicAliasPtr) *b.BasicAliasPtr
)

// go-conv:generate
var (
	ArrayAliToArrayAli    func(*a.ArrayAli) *b.ArrayAli
	ArrayAliToArrayAliPtr func(*a.ArrayAli) *b.ArrayAliPtr
	ArrayAliToSliceAli    func(*a.ArrayAli) *b.SliceAli
	ArrayAliToSliceAliPtr func(*a.ArrayAli) *b.SliceAliPtr

	ArrayAliPtrToArrayAli    func(*a.ArrayAliPtr) *b.ArrayAli
	ArrayAliPtrToArrayAliPtr func(*a.ArrayAliPtr) *b.ArrayAliPtr
	ArrayAliPtrToSliceAli    func(*a.ArrayAliPtr) *b.SliceAli
	ArrayAliPtrToSliceAliPtr func(*a.ArrayAliPtr) *b.SliceAliPtr

	SliceAliToSliceAli    func(*a.SliceAli) *b.SliceAli
	SliceAliToSliceAliPtr func(*a.SliceAli) *b.SliceAliPtr
	SliceAliToArrayAli    func(*a.SliceAli) *b.ArrayAli
	SliceAliToArrayAliPtr func(*a.SliceAli) *b.ArrayAliPtr

	SliceAliPtrToArrayAli    func(*a.SliceAliPtr) *b.ArrayAli
	SliceAliPtrToArrayAliPtr func(*a.SliceAliPtr) *b.ArrayAliPtr
	SliceAliPtrToSliceAli    func(*a.SliceAliPtr) *b.SliceAli
	SliceAliPtrToSliceAliPtr func(*a.SliceAliPtr) *b.SliceAliPtr

	MapToMap       func(*a.MapAli) *b.MapAli
	MapToMapPtr    func(*a.MapAli) *b.MapAliPtr
	MapPtrToMap    func(ptr *a.MapAliPtr) *b.MapAli
	MapPtrToMapPtr func(ptr *a.MapAliPtr) *b.MapAliPtr
)

// go-conv:generate
// go-conv:copy
var (
	CopyArrayAliToArrayAli    func(*a.ArrayAli) *b.ArrayAli
	CopyArrayAliToArrayAliPtr func(*a.ArrayAli) *b.ArrayAliPtr
	CopyArrayAliToSliceAli    func(*a.ArrayAli) *b.SliceAli
	CopyArrayAliToSliceAliPtr func(*a.ArrayAli) *b.SliceAliPtr

	CopyArrayAliPtrToArrayAli    func(*a.ArrayAliPtr) *b.ArrayAli
	CopyArrayAliPtrToArrayAliPtr func(*a.ArrayAliPtr) *b.ArrayAliPtr
	CopyArrayAliPtrToSliceAli    func(*a.ArrayAliPtr) *b.SliceAli
	CopyArrayAliPtrToSliceAliPtr func(*a.ArrayAliPtr) *b.SliceAliPtr

	CopySliceAliToSliceAli    func(*a.SliceAli) *b.SliceAli
	CopySliceAliToSliceAliPtr func(*a.SliceAli) *b.SliceAliPtr
	CopySliceAliToArrayAli    func(*a.SliceAli) *b.ArrayAli
	CopySliceAliToArrayAliPtr func(*a.SliceAli) *b.ArrayAliPtr

	CopySliceAliPtrToArrayAli    func(*a.SliceAliPtr) *b.ArrayAli
	CopySliceAliPtrToArrayAliPtr func(*a.SliceAliPtr) *b.ArrayAliPtr
	CopySliceAliPtrToSliceAli    func(*a.SliceAliPtr) *b.SliceAli
	CopySliceAliPtrToSliceAliPtr func(*a.SliceAliPtr) *b.SliceAliPtr

	CopyMapToMap       func(*a.MapAli) *b.MapAli
	CopyMapToMapPtr    func(*a.MapAli) *b.MapAliPtr
	CopyMapPtrToMap    func(ptr *a.MapAliPtr) *b.MapAli
	CopyMapPtrToMapPtr func(ptr *a.MapAliPtr) *b.MapAliPtr
)

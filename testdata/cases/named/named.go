package named

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
var (
	BasicToBasic               func(*a.BasicNamed) *b.BasicNamed
	BasicToBasicPointer        func(*a.BasicNamed) *b.BasicNamedPtr
	BasicPointerToBasic        func(ptr *a.BasicNamedPtr) *b.BasicNamed
	BasicPointerToBasicPointer func(ptr *a.BasicNamedPtr) *b.BasicNamedPtr
)

// go-conv:generate
// go-conv:copy
var (
	CopyBasicToBasic               func(*a.BasicNamed) *b.BasicNamed
	CopyBasicToBasicPointer        func(*a.BasicNamed) *b.BasicNamedPtr
	CopyBasicPointerToBasic        func(ptr *a.BasicNamedPtr) *b.BasicNamed
	CopyBasicPointerToBasicPointer func(ptr *a.BasicNamedPtr) *b.BasicNamedPtr
)

// go-conv:generate
var (
	ArrayNToArrayN    func(*a.ArrayN) *b.ArrayN
	ArrayNToArrayNPtr func(*a.ArrayN) *b.ArrayNPtr
	ArrayNToSliceN    func(*a.ArrayN) *b.SliceN
	ArrayNToSliceNPtr func(*a.ArrayN) *b.SliceNPtr

	ArrayNPtrToArrayN    func(*a.ArrayNPtr) *b.ArrayN
	ArrayNPtrToArrayNPtr func(*a.ArrayNPtr) *b.ArrayNPtr
	ArrayNPtrToSliceN    func(*a.ArrayNPtr) *b.SliceN
	ArrayNPtrToSliceNPtr func(*a.ArrayNPtr) *b.SliceNPtr

	SliceNToSliceN    func(*a.SliceN) *b.SliceN
	SliceNToSliceNPtr func(*a.SliceN) *b.SliceNPtr
	SliceNToArrayN    func(*a.SliceN) *b.ArrayN
	SliceNToArrayNPtr func(*a.SliceN) *b.ArrayNPtr

	SliceNPtrToArrayN    func(*a.SliceNPtr) *b.ArrayN
	SliceNPtrToArrayNPtr func(*a.SliceNPtr) *b.ArrayNPtr
	SliceNPtrToSliceN    func(*a.SliceNPtr) *b.SliceN
	SliceNPtrToSliceNPtr func(*a.SliceNPtr) *b.SliceNPtr

	MapNToMapN       func(*a.MapN) *b.MapN
	MapNToMapNPtr    func(*a.MapN) *b.MapNPtr
	MapNPtrToMapN    func(ptr *a.MapNPtr) *b.MapN
	MapNPtrToMapNPtr func(ptr *a.MapNPtr) *b.MapNPtr
)

// go-conv:generate
// go-conv:copy
var (
	CopyArrayNToArrayN    func(*a.ArrayN) *b.ArrayN
	CopyArrayNToArrayNPtr func(*a.ArrayN) *b.ArrayNPtr
	CopyArrayNToSliceN    func(*a.ArrayN) *b.SliceN
	CopyArrayNToSliceNPtr func(*a.ArrayN) *b.SliceNPtr

	CopyArrayNPtrToArrayN    func(*a.ArrayNPtr) *b.ArrayN
	CopyArrayNPtrToArrayNPtr func(*a.ArrayNPtr) *b.ArrayNPtr
	CopyArrayNPtrToSliceN    func(*a.ArrayNPtr) *b.SliceN
	CopyArrayNPtrToSliceNPtr func(*a.ArrayNPtr) *b.SliceNPtr

	CopySliceNToSliceN    func(*a.SliceN) *b.SliceN
	CopySliceNToSliceNPtr func(*a.SliceN) *b.SliceNPtr
	CopySliceNToArrayN    func(*a.SliceN) *b.ArrayN
	CopySliceNToArrayNPtr func(*a.SliceN) *b.ArrayNPtr

	CopySliceNPtrToArrayN    func(*a.SliceNPtr) *b.ArrayN
	CopySliceNPtrToArrayNPtr func(*a.SliceNPtr) *b.ArrayNPtr
	CopySliceNPtrToSliceN    func(*a.SliceNPtr) *b.SliceN
	CopySliceNPtrToSliceNPtr func(*a.SliceNPtr) *b.SliceNPtr

	CopyMapNToMapN       func(*a.MapN) *b.MapN
	CopyMapNToMapNPtr    func(*a.MapN) *b.MapNPtr
	CopyMapNPtrToMapN    func(ptr *a.MapNPtr) *b.MapN
	CopyMapNPtrToMapNPtr func(ptr *a.MapNPtr) *b.MapNPtr
)

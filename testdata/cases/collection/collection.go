package collection

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

// go-conv:generate
//var ArrayPtrToSlicePtrTest func(*a.ArrayPtr) *b.SlicePtr

// go-conv:generate
var (
	ArrayToArray    func(*a.Array) *b.Array
	ArrayToArrayPtr func(*a.Array) *b.ArrayPtr
	ArrayToSlice    func(*a.Array) *b.Slice
	ArrayToSlicePtr func(*a.Array) *b.SlicePtr

	ArrayPtrToArray    func(*a.ArrayPtr) *b.Array
	ArrayPtrToArrayPtr func(*a.ArrayPtr) *b.ArrayPtr
	ArrayPtrToSlice    func(*a.ArrayPtr) *b.Slice
	ArrayPtrToSlicePtr func(*a.ArrayPtr) *b.SlicePtr

	SliceToSlice    func(*a.Slice) *b.Slice
	SliceToSlicePtr func(*a.Slice) *b.SlicePtr
	SliceToArray    func(*a.Slice) *b.Array
	SliceToArrayPtr func(*a.Slice) *b.ArrayPtr

	SlicePtrToArray    func(*a.SlicePtr) *b.Array
	SlicePtrToArrayPtr func(*a.SlicePtr) *b.ArrayPtr
	SlicePtrToSlice    func(*a.SlicePtr) *b.Slice
	SlicePtrToSlicePtr func(*a.SlicePtr) *b.SlicePtr

	MapToMap       func(*a.Map) *b.Map
	MapToMapPtr    func(*a.Map) *b.MapPtr
	MapPtrToMap    func(ptr *a.MapPtr) *b.Map
	MapPtrToMapPtr func(ptr *a.MapPtr) *b.MapPtr
)

// go-conv:generate
// go-conv:copy
var (
	CopyArrayToArray    func(*a.Array) *b.Array
	CopyArrayToArrayPtr func(*a.Array) *b.ArrayPtr
	CopyArrayToSlice    func(*a.Array) *b.Slice
	CopyArrayToSlicePtr func(*a.Array) *b.SlicePtr

	CopyArrayPtrToArray    func(*a.ArrayPtr) *b.Array
	CopyArrayPtrToArrayPtr func(*a.ArrayPtr) *b.ArrayPtr
	CopyArrayPtrToSlice    func(*a.ArrayPtr) *b.Slice
	CopyArrayPtrToSlicePtr func(*a.ArrayPtr) *b.SlicePtr

	CopySliceToSlice    func(*a.Slice) *b.Slice
	CopySliceToSlicePtr func(*a.Slice) *b.SlicePtr
	CopySliceToArray    func(*a.Slice) *b.Array
	CopySliceToArrayPtr func(*a.Slice) *b.ArrayPtr

	CopySlicePtrToArray    func(*a.SlicePtr) *b.Array
	CopySlicePtrToArrayPtr func(*a.SlicePtr) *b.ArrayPtr
	CopySlicePtrToSlice    func(*a.SlicePtr) *b.Slice
	CopySlicePtrToSlicePtr func(*a.SlicePtr) *b.SlicePtr

	CopyMapToMap       func(*a.Map) *b.Map
	CopyMapToMapPtr    func(*a.Map) *b.MapPtr
	CopyMapPtrToMap    func(ptr *a.MapPtr) *b.Map
	CopyMapPtrToMapPtr func(ptr *a.MapPtr) *b.MapPtr
)

// Code generated by github.com/ycl2018/go-conv DO NOT EDIT.

package collection

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

func CopyPtrAArrayPtrToPtrBArray(src *a.ArrayPtr) (dst *b.Array) {
	if src != nil {
		dst = new(b.Array)
		for i := 0; i < 6; i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrAArrayPtrToPtrBArrayPtr(src *a.ArrayPtr) (dst *b.ArrayPtr) {
	if src != nil {
		dst = new(b.ArrayPtr)
		*dst = b.ArrayPtr(*src)
	}
	return
}
func CopyPtrAArrayPtrToPtrBSlice(src *a.ArrayPtr) (dst *b.Slice) {
	if src != nil {
		dst = new(b.Slice)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrAArrayPtrToPtrBSlicePtr(src *a.ArrayPtr) (dst *b.SlicePtr) {
	if src != nil {
		dst = new(b.SlicePtr)
		if src.Name != nil {
			dst.Name = new([]string)
			*dst.Name = make([]string, len((*src.Name)))
			for i := 0; i < len((*src.Name)); i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func CopyPtrAArrayToPtrBArray(src *a.Array) (dst *b.Array) {
	if src != nil {
		dst = new(b.Array)
		*dst = b.Array(*src)
	}
	return
}
func CopyPtrAArrayToPtrBArrayPtr(src *a.Array) (dst *b.ArrayPtr) {
	if src != nil {
		dst = new(b.ArrayPtr)
		dst.Name = new([6]string)
		for i := 0; i < 6; i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrAArrayToPtrBSlice(src *a.Array) (dst *b.Slice) {
	if src != nil {
		dst = new(b.Slice)
		dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrAArrayToPtrBSlicePtr(src *a.Array) (dst *b.SlicePtr) {
	if src != nil {
		dst = new(b.SlicePtr)
		dst.Name = new([]string)
		*dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrAMapPtrToPtrBMap(src *a.MapPtr) (dst *b.Map) {
	if src != nil {
		dst = new(b.Map)
		if len((*src.Name)) > 0 {
			dst.Name = make(map[string]string, len((*src.Name)))
			for k, v := range *src.Name {
				var tmpK string
				var tmpV string
				tmpK = k
				tmpV = v
				dst.Name[tmpK] = tmpV
			}
		}
	}
	return
}
func CopyPtrAMapPtrToPtrBMapPtr(src *a.MapPtr) (dst *b.MapPtr) {
	if src != nil {
		dst = new(b.MapPtr)
		*dst = b.MapPtr(*src)
	}
	return
}
func CopyPtrAMapToPtrBMap(src *a.Map) (dst *b.Map) {
	if src != nil {
		dst = new(b.Map)
		*dst = b.Map(*src)
	}
	return
}
func CopyPtrAMapToPtrBMapPtr(src *a.Map) (dst *b.MapPtr) {
	if src != nil {
		dst = new(b.MapPtr)
		dst.Name = new(map[string]string)
		if len(src.Name) > 0 {
			*dst.Name = make(map[string]string, len(src.Name))
			for k, v := range src.Name {
				var tmpK string
				var tmpV string
				tmpK = k
				tmpV = v
				(*dst.Name)[tmpK] = tmpV
			}
		}
	}
	return
}
func CopyPtrASlicePtrToPtrBArray(src *a.SlicePtr) (dst *b.Array) {
	if src != nil {
		dst = new(b.Array)
		for i := 0; i < 6 && i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrASlicePtrToPtrBArrayPtr(src *a.SlicePtr) (dst *b.ArrayPtr) {
	if src != nil {
		dst = new(b.ArrayPtr)
		if src.Name != nil {
			dst.Name = new([6]string)
			for i := 0; i < 6 && i < len((*src.Name)); i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func CopyPtrASlicePtrToPtrBSlice(src *a.SlicePtr) (dst *b.Slice) {
	if src != nil {
		dst = new(b.Slice)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrASlicePtrToPtrBSlicePtr(src *a.SlicePtr) (dst *b.SlicePtr) {
	if src != nil {
		dst = new(b.SlicePtr)
		*dst = b.SlicePtr(*src)
	}
	return
}
func CopyPtrASliceToPtrBArray(src *a.Slice) (dst *b.Array) {
	if src != nil {
		dst = new(b.Array)
		for i := 0; i < 6 && i < len(src.Name); i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrASliceToPtrBArrayPtr(src *a.Slice) (dst *b.ArrayPtr) {
	if src != nil {
		dst = new(b.ArrayPtr)
		dst.Name = new([6]string)
		for i := 0; i < 6 && i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrASliceToPtrBSlice(src *a.Slice) (dst *b.Slice) {
	if src != nil {
		dst = new(b.Slice)
		*dst = b.Slice(*src)
	}
	return
}
func CopyPtrASliceToPtrBSlicePtr(src *a.Slice) (dst *b.SlicePtr) {
	if src != nil {
		dst = new(b.SlicePtr)
		dst.Name = new([]string)
		*dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func PtrAArrayPtrToPtrBArray(src *a.ArrayPtr) (dst *b.Array) {
	if src != nil {
		dst = new(b.Array)
		for i := 0; i < 6; i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrAArrayPtrToPtrBArrayPtr(src *a.ArrayPtr) (dst *b.ArrayPtr) {
	dst = (*b.ArrayPtr)(src)
	return
}
func PtrAArrayPtrToPtrBSlice(src *a.ArrayPtr) (dst *b.Slice) {
	if src != nil {
		dst = new(b.Slice)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrAArrayPtrToPtrBSlicePtr(src *a.ArrayPtr) (dst *b.SlicePtr) {
	if src != nil {
		dst = new(b.SlicePtr)
		if src.Name != nil {
			dst.Name = new([]string)
			*dst.Name = make([]string, len((*src.Name)))
			for i := 0; i < len((*src.Name)); i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func PtrAArrayToPtrBArray(src *a.Array) (dst *b.Array) {
	dst = (*b.Array)(src)
	return
}
func PtrAArrayToPtrBArrayPtr(src *a.Array) (dst *b.ArrayPtr) {
	if src != nil {
		dst = new(b.ArrayPtr)
		dst.Name = new([6]string)
		*dst.Name = src.Name
	}
	return
}
func PtrAArrayToPtrBSlice(src *a.Array) (dst *b.Slice) {
	if src != nil {
		dst = new(b.Slice)
		dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func PtrAArrayToPtrBSlicePtr(src *a.Array) (dst *b.SlicePtr) {
	if src != nil {
		dst = new(b.SlicePtr)
		dst.Name = new([]string)
		*dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func PtrAMapPtrToPtrBMap(src *a.MapPtr) (dst *b.Map) {
	if src != nil {
		dst = new(b.Map)
		if len((*src.Name)) > 0 {
			dst.Name = make(map[string]string, len((*src.Name)))
			for k, v := range *src.Name {
				var tmpK string
				var tmpV string
				tmpK = k
				tmpV = v
				dst.Name[tmpK] = tmpV
			}
		}
	}
	return
}
func PtrAMapPtrToPtrBMapPtr(src *a.MapPtr) (dst *b.MapPtr) {
	dst = (*b.MapPtr)(src)
	return
}
func PtrAMapToPtrBMap(src *a.Map) (dst *b.Map) {
	dst = (*b.Map)(src)
	return
}
func PtrAMapToPtrBMapPtr(src *a.Map) (dst *b.MapPtr) {
	if src != nil {
		dst = new(b.MapPtr)
		dst.Name = new(map[string]string)
		*dst.Name = src.Name
	}
	return
}
func PtrASlicePtrToPtrBArray(src *a.SlicePtr) (dst *b.Array) {
	if src != nil {
		dst = new(b.Array)
		for i := 0; i < 6 && i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrASlicePtrToPtrBArrayPtr(src *a.SlicePtr) (dst *b.ArrayPtr) {
	if src != nil {
		dst = new(b.ArrayPtr)
		if src.Name != nil {
			dst.Name = new([6]string)
			*dst.Name = [6]string(*src.Name)
		}
	}
	return
}
func PtrASlicePtrToPtrBSlice(src *a.SlicePtr) (dst *b.Slice) {
	if src != nil {
		dst = new(b.Slice)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrASlicePtrToPtrBSlicePtr(src *a.SlicePtr) (dst *b.SlicePtr) {
	dst = (*b.SlicePtr)(src)
	return
}
func PtrASliceToPtrBArray(src *a.Slice) (dst *b.Array) {
	if src != nil {
		dst = new(b.Array)
		dst.Name = [6]string(src.Name)
	}
	return
}
func PtrASliceToPtrBArrayPtr(src *a.Slice) (dst *b.ArrayPtr) {
	if src != nil {
		dst = new(b.ArrayPtr)
		dst.Name = (*[6]string)(src.Name)
	}
	return
}
func PtrASliceToPtrBSlice(src *a.Slice) (dst *b.Slice) {
	dst = (*b.Slice)(src)
	return
}
func PtrASliceToPtrBSlicePtr(src *a.Slice) (dst *b.SlicePtr) {
	if src != nil {
		dst = new(b.SlicePtr)
		dst.Name = new([]string)
		*dst.Name = src.Name
	}
	return
}
func init() {
	ArrayPtrToArray = PtrAArrayPtrToPtrBArray
	ArrayPtrToArrayPtr = PtrAArrayPtrToPtrBArrayPtr
	ArrayPtrToSlice = PtrAArrayPtrToPtrBSlice
	ArrayPtrToSlicePtr = PtrAArrayPtrToPtrBSlicePtr
	ArrayToArray = PtrAArrayToPtrBArray
	ArrayToArrayPtr = PtrAArrayToPtrBArrayPtr
	ArrayToSlice = PtrAArrayToPtrBSlice
	ArrayToSlicePtr = PtrAArrayToPtrBSlicePtr
	CopyArrayPtrToArray = CopyPtrAArrayPtrToPtrBArray
	CopyArrayPtrToArrayPtr = CopyPtrAArrayPtrToPtrBArrayPtr
	CopyArrayPtrToSlice = CopyPtrAArrayPtrToPtrBSlice
	CopyArrayPtrToSlicePtr = CopyPtrAArrayPtrToPtrBSlicePtr
	CopyArrayToArray = CopyPtrAArrayToPtrBArray
	CopyArrayToArrayPtr = CopyPtrAArrayToPtrBArrayPtr
	CopyArrayToSlice = CopyPtrAArrayToPtrBSlice
	CopyArrayToSlicePtr = CopyPtrAArrayToPtrBSlicePtr
	CopyMapPtrToMap = CopyPtrAMapPtrToPtrBMap
	CopyMapPtrToMapPtr = CopyPtrAMapPtrToPtrBMapPtr
	CopyMapToMap = CopyPtrAMapToPtrBMap
	CopyMapToMapPtr = CopyPtrAMapToPtrBMapPtr
	CopySlicePtrToArray = CopyPtrASlicePtrToPtrBArray
	CopySlicePtrToArrayPtr = CopyPtrASlicePtrToPtrBArrayPtr
	CopySlicePtrToSlice = CopyPtrASlicePtrToPtrBSlice
	CopySlicePtrToSlicePtr = CopyPtrASlicePtrToPtrBSlicePtr
	CopySliceToArray = CopyPtrASliceToPtrBArray
	CopySliceToArrayPtr = CopyPtrASliceToPtrBArrayPtr
	CopySliceToSlice = CopyPtrASliceToPtrBSlice
	CopySliceToSlicePtr = CopyPtrASliceToPtrBSlicePtr
	MapPtrToMap = PtrAMapPtrToPtrBMap
	MapPtrToMapPtr = PtrAMapPtrToPtrBMapPtr
	MapToMap = PtrAMapToPtrBMap
	MapToMapPtr = PtrAMapToPtrBMapPtr
	SlicePtrToArray = PtrASlicePtrToPtrBArray
	SlicePtrToArrayPtr = PtrASlicePtrToPtrBArrayPtr
	SlicePtrToSlice = PtrASlicePtrToPtrBSlice
	SlicePtrToSlicePtr = PtrASlicePtrToPtrBSlicePtr
	SliceToArray = PtrASliceToPtrBArray
	SliceToArrayPtr = PtrASliceToPtrBArrayPtr
	SliceToSlice = PtrASliceToPtrBSlice
	SliceToSlicePtr = PtrASliceToPtrBSlicePtr
}

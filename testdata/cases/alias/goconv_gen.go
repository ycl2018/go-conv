// Code generated by github.com/ycl2018/go-conv DO NOT EDIT.

package alias

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
	"unsafe"
)

func CopyPtrAArrayAliPtrToPtrBArrayAli(src *a.ArrayAliPtr) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		for i := 0; i < 6; i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrAArrayAliPtrToPtrBArrayAliPtr(src *a.ArrayAliPtr) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		if src.Name != nil {
			dst.Name = new(b.AliasArray)
			for i := 0; i < 6; i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func CopyPtrAArrayAliPtrToPtrBSliceAli(src *a.ArrayAliPtr) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrAArrayAliPtrToPtrBSliceAliPtr(src *a.ArrayAliPtr) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		if src.Name != nil {
			dst.Name = new(b.AliasSlice)
			*dst.Name = make([]string, len((*src.Name)))
			for i := 0; i < len((*src.Name)); i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func CopyPtrAArrayAliToPtrBArrayAli(src *a.ArrayAli) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		for i := 0; i < 6; i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrAArrayAliToPtrBArrayAliPtr(src *a.ArrayAli) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		dst.Name = new(b.AliasArray)
		for i := 0; i < 6; i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrAArrayAliToPtrBSliceAli(src *a.ArrayAli) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrAArrayAliToPtrBSliceAliPtr(src *a.ArrayAli) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		dst.Name = new(b.AliasSlice)
		*dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrABasicAliasPtrToPtrBBasicAlias(src *a.BasicAliasPtr) (dst *b.BasicAlias) {
	if src != nil {
		dst = new(b.BasicAlias)
		dst.Bool = *src.Bool
		dst.Int = *src.Int
		dst.Int8 = *src.Int8
		dst.Int16 = *src.Int16
		dst.Int32 = *src.Int32
		dst.Int64 = *src.Int64
		dst.Uint = *src.Uint
		dst.Uint8 = *src.Uint8
		dst.Uint16 = *src.Uint16
		dst.Uint32 = *src.Uint32
		dst.Uint64 = *src.Uint64
		dst.Uintptr = *src.Uintptr
		dst.Float32 = *src.Float32
		dst.Float64 = *src.Float64
		dst.Complex64 = *src.Complex64
		dst.Complex128 = *src.Complex128
		dst.String = *src.String
		dst.UnsafePointer = *src.UnsafePointer
		dst.Byte = *src.Byte
		dst.Rune = *src.Rune
	}
	return
}
func CopyPtrABasicAliasPtrToPtrBBasicAliasPtr(src *a.BasicAliasPtr) (dst *b.BasicAliasPtr) {
	if src != nil {
		dst = new(b.BasicAliasPtr)
		*dst = b.BasicAliasPtr(*src)
	}
	return
}
func CopyPtrABasicAliasToPtrBBasicAlias(src *a.BasicAlias) (dst *b.BasicAlias) {
	if src != nil {
		dst = new(b.BasicAlias)
		*dst = b.BasicAlias(*src)
	}
	return
}
func CopyPtrABasicAliasToPtrBBasicAliasPtr(src *a.BasicAlias) (dst *b.BasicAliasPtr) {
	if src != nil {
		dst = new(b.BasicAliasPtr)
		dst.Bool = new(bool)
		*dst.Bool = src.Bool
		dst.Int = new(int)
		*dst.Int = src.Int
		dst.Int8 = new(int8)
		*dst.Int8 = src.Int8
		dst.Int16 = new(int16)
		*dst.Int16 = src.Int16
		dst.Int32 = new(int32)
		*dst.Int32 = src.Int32
		dst.Int64 = new(int64)
		*dst.Int64 = src.Int64
		dst.Uint = new(uint)
		*dst.Uint = src.Uint
		dst.Uint8 = new(uint8)
		*dst.Uint8 = src.Uint8
		dst.Uint16 = new(uint16)
		*dst.Uint16 = src.Uint16
		dst.Uint32 = new(uint32)
		*dst.Uint32 = src.Uint32
		dst.Uint64 = new(uint64)
		*dst.Uint64 = src.Uint64
		dst.Uintptr = new(uintptr)
		*dst.Uintptr = src.Uintptr
		dst.Float32 = new(float32)
		*dst.Float32 = src.Float32
		dst.Float64 = new(float64)
		*dst.Float64 = src.Float64
		dst.Complex64 = new(complex64)
		*dst.Complex64 = src.Complex64
		dst.Complex128 = new(complex128)
		*dst.Complex128 = src.Complex128
		dst.String = new(string)
		*dst.String = src.String
		dst.UnsafePointer = new(unsafe.Pointer)
		*dst.UnsafePointer = src.UnsafePointer
		dst.Byte = new(byte)
		*dst.Byte = src.Byte
		dst.Rune = new(rune)
		*dst.Rune = src.Rune
	}
	return
}
func CopyPtrAMapAliPtrToPtrBMapAli(src *a.MapAliPtr) (dst *b.MapAli) {
	if src != nil {
		dst = new(b.MapAli)
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
func CopyPtrAMapAliPtrToPtrBMapAliPtr(src *a.MapAliPtr) (dst *b.MapAliPtr) {
	if src != nil {
		dst = new(b.MapAliPtr)
		if src.Name != nil {
			dst.Name = new(b.AliasMap)
			if len((*src.Name)) > 0 {
				*dst.Name = make(map[string]string, len((*src.Name)))
				for k, v := range *src.Name {
					var tmpK string
					var tmpV string
					tmpK = k
					tmpV = v
					(*dst.Name)[tmpK] = tmpV
				}
			}
		}
	}
	return
}
func CopyPtrAMapAliToPtrBMapAli(src *a.MapAli) (dst *b.MapAli) {
	if src != nil {
		dst = new(b.MapAli)
		if len(src.Name) > 0 {
			dst.Name = make(map[string]string, len(src.Name))
			for k, v := range src.Name {
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
func CopyPtrAMapAliToPtrBMapAliPtr(src *a.MapAli) (dst *b.MapAliPtr) {
	if src != nil {
		dst = new(b.MapAliPtr)
		dst.Name = new(b.AliasMap)
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
func CopyPtrASliceAliPtrToPtrBArrayAli(src *a.SliceAliPtr) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		for i := 0; i < 6 && i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrASliceAliPtrToPtrBArrayAliPtr(src *a.SliceAliPtr) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		if src.Name != nil {
			dst.Name = new(b.AliasArray)
			for i := 0; i < 6 && i < len((*src.Name)); i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func CopyPtrASliceAliPtrToPtrBSliceAli(src *a.SliceAliPtr) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func CopyPtrASliceAliPtrToPtrBSliceAliPtr(src *a.SliceAliPtr) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		if src.Name != nil {
			dst.Name = new(b.AliasSlice)
			*dst.Name = make([]string, len((*src.Name)))
			for i := 0; i < len((*src.Name)); i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func CopyPtrASliceAliToPtrBArrayAli(src *a.SliceAli) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		for i := 0; i < 6 && i < len(src.Name); i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrASliceAliToPtrBArrayAliPtr(src *a.SliceAli) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		dst.Name = new(b.AliasArray)
		for i := 0; i < 6 && i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrASliceAliToPtrBSliceAli(src *a.SliceAli) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func CopyPtrASliceAliToPtrBSliceAliPtr(src *a.SliceAli) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		dst.Name = new(b.AliasSlice)
		*dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func PtrAArrayAliPtrToPtrBArrayAli(src *a.ArrayAliPtr) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		for i := 0; i < 6; i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrAArrayAliPtrToPtrBArrayAliPtr(src *a.ArrayAliPtr) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		dst.Name = (*b.AliasArray)(src.Name)
	}
	return
}
func PtrAArrayAliPtrToPtrBSliceAli(src *a.ArrayAliPtr) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrAArrayAliPtrToPtrBSliceAliPtr(src *a.ArrayAliPtr) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		if src.Name != nil {
			dst.Name = new(b.AliasSlice)
			*dst.Name = make([]string, len((*src.Name)))
			for i := 0; i < len((*src.Name)); i++ {
				(*dst.Name)[i] = (*src.Name)[i]
			}
		}
	}
	return
}
func PtrAArrayAliToPtrBArrayAli(src *a.ArrayAli) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		dst.Name = b.AliasArray(src.Name)
	}
	return
}
func PtrAArrayAliToPtrBArrayAliPtr(src *a.ArrayAli) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		dst.Name = new(b.AliasArray)
		*dst.Name = b.AliasArray(src.Name)
	}
	return
}
func PtrAArrayAliToPtrBSliceAli(src *a.ArrayAli) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			dst.Name[i] = src.Name[i]
		}
	}
	return
}
func PtrAArrayAliToPtrBSliceAliPtr(src *a.ArrayAli) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		dst.Name = new(b.AliasSlice)
		*dst.Name = make([]string, len(src.Name))
		for i := 0; i < len(src.Name); i++ {
			(*dst.Name)[i] = src.Name[i]
		}
	}
	return
}
func PtrABasicAliasPtrToPtrBBasicAlias(src *a.BasicAliasPtr) (dst *b.BasicAlias) {
	if src != nil {
		dst = new(b.BasicAlias)
		dst.Bool = *src.Bool
		dst.Int = *src.Int
		dst.Int8 = *src.Int8
		dst.Int16 = *src.Int16
		dst.Int32 = *src.Int32
		dst.Int64 = *src.Int64
		dst.Uint = *src.Uint
		dst.Uint8 = *src.Uint8
		dst.Uint16 = *src.Uint16
		dst.Uint32 = *src.Uint32
		dst.Uint64 = *src.Uint64
		dst.Uintptr = *src.Uintptr
		dst.Float32 = *src.Float32
		dst.Float64 = *src.Float64
		dst.Complex64 = *src.Complex64
		dst.Complex128 = *src.Complex128
		dst.String = *src.String
		dst.UnsafePointer = unsafe.Pointer(src.UnsafePointer)
		dst.Byte = *src.Byte
		dst.Rune = *src.Rune
	}
	return
}
func PtrABasicAliasPtrToPtrBBasicAliasPtr(src *a.BasicAliasPtr) (dst *b.BasicAliasPtr) {
	dst = (*b.BasicAliasPtr)(src)
	return
}
func PtrABasicAliasToPtrBBasicAlias(src *a.BasicAlias) (dst *b.BasicAlias) {
	dst = (*b.BasicAlias)(src)
	return
}
func PtrABasicAliasToPtrBBasicAliasPtr(src *a.BasicAlias) (dst *b.BasicAliasPtr) {
	if src != nil {
		dst = new(b.BasicAliasPtr)
		dst.Bool = new(bool)
		*dst.Bool = src.Bool
		dst.Int = new(int)
		*dst.Int = src.Int
		dst.Int8 = new(int8)
		*dst.Int8 = src.Int8
		dst.Int16 = new(int16)
		*dst.Int16 = src.Int16
		dst.Int32 = new(int32)
		*dst.Int32 = src.Int32
		dst.Int64 = new(int64)
		*dst.Int64 = src.Int64
		dst.Uint = new(uint)
		*dst.Uint = src.Uint
		dst.Uint8 = new(uint8)
		*dst.Uint8 = src.Uint8
		dst.Uint16 = new(uint16)
		*dst.Uint16 = src.Uint16
		dst.Uint32 = new(uint32)
		*dst.Uint32 = src.Uint32
		dst.Uint64 = new(uint64)
		*dst.Uint64 = src.Uint64
		dst.Uintptr = new(uintptr)
		*dst.Uintptr = src.Uintptr
		dst.Float32 = new(float32)
		*dst.Float32 = src.Float32
		dst.Float64 = new(float64)
		*dst.Float64 = src.Float64
		dst.Complex64 = new(complex64)
		*dst.Complex64 = src.Complex64
		dst.Complex128 = new(complex128)
		*dst.Complex128 = src.Complex128
		dst.String = new(string)
		*dst.String = src.String
		dst.UnsafePointer = (*unsafe.Pointer)(src.UnsafePointer)
		dst.Byte = new(byte)
		*dst.Byte = src.Byte
		dst.Rune = new(rune)
		*dst.Rune = src.Rune
	}
	return
}
func PtrAMapAliPtrToPtrBMapAli(src *a.MapAliPtr) (dst *b.MapAli) {
	if src != nil {
		dst = new(b.MapAli)
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
func PtrAMapAliPtrToPtrBMapAliPtr(src *a.MapAliPtr) (dst *b.MapAliPtr) {
	if src != nil {
		dst = new(b.MapAliPtr)
		dst.Name = (*b.AliasMap)(src.Name)
	}
	return
}
func PtrAMapAliToPtrBMapAli(src *a.MapAli) (dst *b.MapAli) {
	if src != nil {
		dst = new(b.MapAli)
		dst.Name = b.AliasMap(src.Name)
	}
	return
}
func PtrAMapAliToPtrBMapAliPtr(src *a.MapAli) (dst *b.MapAliPtr) {
	if src != nil {
		dst = new(b.MapAliPtr)
		dst.Name = new(b.AliasMap)
		*dst.Name = b.AliasMap(src.Name)
	}
	return
}
func PtrASliceAliPtrToPtrBArrayAli(src *a.SliceAliPtr) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		for i := 0; i < 6 && i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrASliceAliPtrToPtrBArrayAliPtr(src *a.SliceAliPtr) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		if src.Name != nil {
			dst.Name = new(b.AliasArray)
			*dst.Name = b.AliasArray(*src.Name)
		}
	}
	return
}
func PtrASliceAliPtrToPtrBSliceAli(src *a.SliceAliPtr) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = make([]string, len((*src.Name)))
		for i := 0; i < len((*src.Name)); i++ {
			dst.Name[i] = (*src.Name)[i]
		}
	}
	return
}
func PtrASliceAliPtrToPtrBSliceAliPtr(src *a.SliceAliPtr) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		dst.Name = (*b.AliasSlice)(src.Name)
	}
	return
}
func PtrASliceAliToPtrBArrayAli(src *a.SliceAli) (dst *b.ArrayAli) {
	if src != nil {
		dst = new(b.ArrayAli)
		dst.Name = b.AliasArray(src.Name)
	}
	return
}
func PtrASliceAliToPtrBArrayAliPtr(src *a.SliceAli) (dst *b.ArrayAliPtr) {
	if src != nil {
		dst = new(b.ArrayAliPtr)
		dst.Name = (*b.AliasArray)(src.Name)
	}
	return
}
func PtrASliceAliToPtrBSliceAli(src *a.SliceAli) (dst *b.SliceAli) {
	if src != nil {
		dst = new(b.SliceAli)
		dst.Name = b.AliasSlice(src.Name)
	}
	return
}
func PtrASliceAliToPtrBSliceAliPtr(src *a.SliceAli) (dst *b.SliceAliPtr) {
	if src != nil {
		dst = new(b.SliceAliPtr)
		dst.Name = new(b.AliasSlice)
		*dst.Name = b.AliasSlice(src.Name)
	}
	return
}
func init() {
	ArrayAliPtrToArrayAli = PtrAArrayAliPtrToPtrBArrayAli
	ArrayAliPtrToArrayAliPtr = PtrAArrayAliPtrToPtrBArrayAliPtr
	ArrayAliPtrToSliceAli = PtrAArrayAliPtrToPtrBSliceAli
	ArrayAliPtrToSliceAliPtr = PtrAArrayAliPtrToPtrBSliceAliPtr
	ArrayAliToArrayAli = PtrAArrayAliToPtrBArrayAli
	ArrayAliToArrayAliPtr = PtrAArrayAliToPtrBArrayAliPtr
	ArrayAliToSliceAli = PtrAArrayAliToPtrBSliceAli
	ArrayAliToSliceAliPtr = PtrAArrayAliToPtrBSliceAliPtr
	BasicPointerToBasic = PtrABasicAliasPtrToPtrBBasicAlias
	BasicPointerToBasicPointer = PtrABasicAliasPtrToPtrBBasicAliasPtr
	BasicToBasic = PtrABasicAliasToPtrBBasicAlias
	BasicToBasicPointer = PtrABasicAliasToPtrBBasicAliasPtr
	CopyArrayAliPtrToArrayAli = CopyPtrAArrayAliPtrToPtrBArrayAli
	CopyArrayAliPtrToArrayAliPtr = CopyPtrAArrayAliPtrToPtrBArrayAliPtr
	CopyArrayAliPtrToSliceAli = CopyPtrAArrayAliPtrToPtrBSliceAli
	CopyArrayAliPtrToSliceAliPtr = CopyPtrAArrayAliPtrToPtrBSliceAliPtr
	CopyArrayAliToArrayAli = CopyPtrAArrayAliToPtrBArrayAli
	CopyArrayAliToArrayAliPtr = CopyPtrAArrayAliToPtrBArrayAliPtr
	CopyArrayAliToSliceAli = CopyPtrAArrayAliToPtrBSliceAli
	CopyArrayAliToSliceAliPtr = CopyPtrAArrayAliToPtrBSliceAliPtr
	CopyBasicPointerToBasic = CopyPtrABasicAliasPtrToPtrBBasicAlias
	CopyBasicPointerToBasicPointer = CopyPtrABasicAliasPtrToPtrBBasicAliasPtr
	CopyBasicToBasic = CopyPtrABasicAliasToPtrBBasicAlias
	CopyBasicToBasicPointer = CopyPtrABasicAliasToPtrBBasicAliasPtr
	CopyMapPtrToMap = CopyPtrAMapAliPtrToPtrBMapAli
	CopyMapPtrToMapPtr = CopyPtrAMapAliPtrToPtrBMapAliPtr
	CopyMapToMap = CopyPtrAMapAliToPtrBMapAli
	CopyMapToMapPtr = CopyPtrAMapAliToPtrBMapAliPtr
	CopySliceAliPtrToArrayAli = CopyPtrASliceAliPtrToPtrBArrayAli
	CopySliceAliPtrToArrayAliPtr = CopyPtrASliceAliPtrToPtrBArrayAliPtr
	CopySliceAliPtrToSliceAli = CopyPtrASliceAliPtrToPtrBSliceAli
	CopySliceAliPtrToSliceAliPtr = CopyPtrASliceAliPtrToPtrBSliceAliPtr
	CopySliceAliToArrayAli = CopyPtrASliceAliToPtrBArrayAli
	CopySliceAliToArrayAliPtr = CopyPtrASliceAliToPtrBArrayAliPtr
	CopySliceAliToSliceAli = CopyPtrASliceAliToPtrBSliceAli
	CopySliceAliToSliceAliPtr = CopyPtrASliceAliToPtrBSliceAliPtr
	MapPtrToMap = PtrAMapAliPtrToPtrBMapAli
	MapPtrToMapPtr = PtrAMapAliPtrToPtrBMapAliPtr
	MapToMap = PtrAMapAliToPtrBMapAli
	MapToMapPtr = PtrAMapAliToPtrBMapAliPtr
	SliceAliPtrToArrayAli = PtrASliceAliPtrToPtrBArrayAli
	SliceAliPtrToArrayAliPtr = PtrASliceAliPtrToPtrBArrayAliPtr
	SliceAliPtrToSliceAli = PtrASliceAliPtrToPtrBSliceAli
	SliceAliPtrToSliceAliPtr = PtrASliceAliPtrToPtrBSliceAliPtr
	SliceAliToArrayAli = PtrASliceAliToPtrBArrayAli
	SliceAliToArrayAliPtr = PtrASliceAliToPtrBArrayAliPtr
	SliceAliToSliceAli = PtrASliceAliToPtrBSliceAli
	SliceAliToSliceAliPtr = PtrASliceAliToPtrBSliceAliPtr
}

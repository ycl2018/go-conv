// Code generated by github.com/ycl2018/go-conv DO NOT EDIT.

package basic

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
	"unsafe"
)

func ABasicPtrToBBasic(src *a.BasicPtr) (dst *b.Basic) {
	if src != nil {
		dst = new(b.Basic)
		if src.Bool != nil {
			dst.Bool = *src.Bool
		}
		if src.Int != nil {
			dst.Int = *src.Int
		}
		if src.Int8 != nil {
			dst.Int8 = *src.Int8
		}
		if src.Int16 != nil {
			dst.Int16 = *src.Int16
		}
		if src.Int32 != nil {
			dst.Int32 = *src.Int32
		}
		if src.Int64 != nil {
			dst.Int64 = *src.Int64
		}
		if src.Uint != nil {
			dst.Uint = *src.Uint
		}
		if src.Uint8 != nil {
			dst.Uint8 = *src.Uint8
		}
		if src.Uint16 != nil {
			dst.Uint16 = *src.Uint16
		}
		if src.Uint32 != nil {
			dst.Uint32 = *src.Uint32
		}
		if src.Uint64 != nil {
			dst.Uint64 = *src.Uint64
		}
		if src.Uintptr != nil {
			dst.Uintptr = *src.Uintptr
		}
		if src.Float32 != nil {
			dst.Float32 = *src.Float32
		}
		if src.Float64 != nil {
			dst.Float64 = *src.Float64
		}
		if src.Complex64 != nil {
			dst.Complex64 = *src.Complex64
		}
		if src.Complex128 != nil {
			dst.Complex128 = *src.Complex128
		}
		if src.String != nil {
			dst.String = *src.String
		}
		dst.UnsafePointer = (unsafe.Pointer)(src.UnsafePointer)
		if src.Byte != nil {
			dst.Byte = *src.Byte
		}
		if src.Rune != nil {
			dst.Rune = *src.Rune
		}
	}
	return
}
func ABasicPtrToBBasicPtr(src *a.BasicPtr) (dst *b.BasicPtr) {
	dst = (*b.BasicPtr)(src)
	return
}
func ABasicToBBasic(src *a.Basic) (dst *b.Basic) {
	dst = (*b.Basic)(src)
	return
}
func ABasicToBBasicPtr(src *a.Basic) (dst *b.BasicPtr) {
	if src != nil {
		dst = new(b.BasicPtr)
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
func CopyABasicPtrToBBasic(src *a.BasicPtr) (dst *b.Basic) {
	if src != nil {
		dst = new(b.Basic)
		if src.Bool != nil {
			dst.Bool = *src.Bool
		}
		if src.Int != nil {
			dst.Int = *src.Int
		}
		if src.Int8 != nil {
			dst.Int8 = *src.Int8
		}
		if src.Int16 != nil {
			dst.Int16 = *src.Int16
		}
		if src.Int32 != nil {
			dst.Int32 = *src.Int32
		}
		if src.Int64 != nil {
			dst.Int64 = *src.Int64
		}
		if src.Uint != nil {
			dst.Uint = *src.Uint
		}
		if src.Uint8 != nil {
			dst.Uint8 = *src.Uint8
		}
		if src.Uint16 != nil {
			dst.Uint16 = *src.Uint16
		}
		if src.Uint32 != nil {
			dst.Uint32 = *src.Uint32
		}
		if src.Uint64 != nil {
			dst.Uint64 = *src.Uint64
		}
		if src.Uintptr != nil {
			dst.Uintptr = *src.Uintptr
		}
		if src.Float32 != nil {
			dst.Float32 = *src.Float32
		}
		if src.Float64 != nil {
			dst.Float64 = *src.Float64
		}
		if src.Complex64 != nil {
			dst.Complex64 = *src.Complex64
		}
		if src.Complex128 != nil {
			dst.Complex128 = *src.Complex128
		}
		if src.String != nil {
			dst.String = *src.String
		}
		dst.UnsafePointer = unsafe.Pointer(src.UnsafePointer)
		if src.Byte != nil {
			dst.Byte = *src.Byte
		}
		if src.Rune != nil {
			dst.Rune = *src.Rune
		}
	}
	return
}
func CopyABasicPtrToBBasicPtr(src *a.BasicPtr) (dst *b.BasicPtr) {
	if src != nil {
		dst = new(b.BasicPtr)
		if src.Bool != nil {
			dst.Bool = new(bool)
			*dst.Bool = *src.Bool
		}
		if src.Int != nil {
			dst.Int = new(int)
			*dst.Int = *src.Int
		}
		if src.Int8 != nil {
			dst.Int8 = new(int8)
			*dst.Int8 = *src.Int8
		}
		if src.Int16 != nil {
			dst.Int16 = new(int16)
			*dst.Int16 = *src.Int16
		}
		if src.Int32 != nil {
			dst.Int32 = new(int32)
			*dst.Int32 = *src.Int32
		}
		if src.Int64 != nil {
			dst.Int64 = new(int64)
			*dst.Int64 = *src.Int64
		}
		if src.Uint != nil {
			dst.Uint = new(uint)
			*dst.Uint = *src.Uint
		}
		if src.Uint8 != nil {
			dst.Uint8 = new(uint8)
			*dst.Uint8 = *src.Uint8
		}
		if src.Uint16 != nil {
			dst.Uint16 = new(uint16)
			*dst.Uint16 = *src.Uint16
		}
		if src.Uint32 != nil {
			dst.Uint32 = new(uint32)
			*dst.Uint32 = *src.Uint32
		}
		if src.Uint64 != nil {
			dst.Uint64 = new(uint64)
			*dst.Uint64 = *src.Uint64
		}
		if src.Uintptr != nil {
			dst.Uintptr = new(uintptr)
			*dst.Uintptr = *src.Uintptr
		}
		if src.Float32 != nil {
			dst.Float32 = new(float32)
			*dst.Float32 = *src.Float32
		}
		if src.Float64 != nil {
			dst.Float64 = new(float64)
			*dst.Float64 = *src.Float64
		}
		if src.Complex64 != nil {
			dst.Complex64 = new(complex64)
			*dst.Complex64 = *src.Complex64
		}
		if src.Complex128 != nil {
			dst.Complex128 = new(complex128)
			*dst.Complex128 = *src.Complex128
		}
		if src.String != nil {
			dst.String = new(string)
			*dst.String = *src.String
		}
		if src.UnsafePointer != nil {
			dst.UnsafePointer = new(unsafe.Pointer)
			*dst.UnsafePointer = *src.UnsafePointer
		}
		if src.Byte != nil {
			dst.Byte = new(byte)
			*dst.Byte = *src.Byte
		}
		if src.Rune != nil {
			dst.Rune = new(rune)
			*dst.Rune = *src.Rune
		}
	}
	return
}
func CopyABasicToBBasic(src *a.Basic) (dst *b.Basic) {
	if src != nil {
		dst = new(b.Basic)
		dst.Bool = src.Bool
		dst.Int = src.Int
		dst.Int8 = src.Int8
		dst.Int16 = src.Int16
		dst.Int32 = src.Int32
		dst.Int64 = src.Int64
		dst.Uint = src.Uint
		dst.Uint8 = src.Uint8
		dst.Uint16 = src.Uint16
		dst.Uint32 = src.Uint32
		dst.Uint64 = src.Uint64
		dst.Uintptr = src.Uintptr
		dst.Float32 = src.Float32
		dst.Float64 = src.Float64
		dst.Complex64 = src.Complex64
		dst.Complex128 = src.Complex128
		dst.String = src.String
		dst.UnsafePointer = src.UnsafePointer
		dst.Byte = src.Byte
		dst.Rune = src.Rune
	}
	return
}
func CopyABasicToBBasicPtr(src *a.Basic) (dst *b.BasicPtr) {
	if src != nil {
		dst = new(b.BasicPtr)
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
func init() {
	BasicPointerToBasic = ABasicPtrToBBasic
	BasicPointerToBasicPointer = ABasicPtrToBBasicPtr
	BasicToBasic = ABasicToBBasic
	BasicToBasicPointer = ABasicToBBasicPtr
	CopyBasicPointerToBasic = CopyABasicPtrToBBasic
	CopyBasicPointerToBasicPointer = CopyABasicPtrToBBasicPtr
	CopyBasicToBasic = CopyABasicToBBasic
	CopyBasicToBasicPointer = CopyABasicToBBasicPtr
}
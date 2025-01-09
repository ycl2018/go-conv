package a

import "unsafe"

type Basic struct {
	Bool          bool
	Int           int
	Int8          int8
	Int16         int16
	Int32         int32
	Int64         int64
	Uint          uint
	Uint8         uint8
	Uint16        uint16
	Uint32        uint32
	Uint64        uint64
	Uintptr       uintptr
	Float32       float32
	Float64       float64
	Complex64     complex64
	Complex128    complex128
	String        string
	UnsafePointer unsafe.Pointer
	Byte          byte
	Rune          rune
}

type BasicPtr struct {
	Bool          *bool
	Int           *int
	Int8          *int8
	Int16         *int16
	Int32         *int32
	Int64         *int64
	Uint          *uint
	Uint8         *uint8
	Uint16        *uint16
	Uint32        *uint32
	Uint64        *uint64
	Uintptr       *uintptr
	Float32       *float32
	Float64       *float64
	Complex64     *complex64
	Complex128    *complex128
	String        *string
	UnsafePointer *unsafe.Pointer
	Byte          *byte
	Rune          *rune
	Ptr2          ***string
}

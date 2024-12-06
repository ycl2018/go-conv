package b

import "unsafe"

type (
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
)

type BasicNamed struct {
	Bool          Bool
	Int           Int
	Int8          Int8
	Int16         Int16
	Int32         Int32
	Int64         Int64
	Uint          Uint
	Uint8         Uint8
	Uint16        Uint16
	Uint32        Uint32
	Uint64        Uint64
	Uintptr       Uintptr
	Float32       Float32
	Float64       Float64
	Complex64     Complex64
	Complex128    Complex128
	String        String
	UnsafePointer UnsafePointer
	Byte          Byte
	Rune          Rune
}

type BasicNamedPtr struct {
	Bool          *Bool
	Int           *Int
	Int8          *Int8
	Int16         *Int16
	Int32         *Int32
	Int64         *Int64
	Uint          *Uint
	Uint8         *Uint8
	Uint16        *Uint16
	Uint32        *Uint32
	Uint64        *Uint64
	Uintptr       *Uintptr
	Float32       *Float32
	Float64       *Float64
	Complex64     *Complex64
	Complex128    *Complex128
	String        *String
	UnsafePointer *UnsafePointer
	Byte          *Byte
	Rune          *Rune
}

type (
	NamedArray [6]string
	NamedSlice []string
	NamedMap   map[string]string
)

type ArrayN struct {
	Name NamedArray
}

type SliceN struct {
	Name NamedSlice
}

type MapN struct {
	Name NamedMap
}

type ArrayNPtr struct {
	Name *NamedArray
}

type SliceNPtr struct {
	Name *NamedSlice
}

type MapNPtr struct {
	Name *NamedMap
}

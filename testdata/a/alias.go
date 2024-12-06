package a

import "unsafe"

type (
	AliasBool          = bool
	AliasInt           = int
	AliasInt8          = int8
	AliasInt16         = int16
	AliasInt32         = int32
	AliasInt64         = int64
	AliasUint          = uint
	AliasUint8         = uint8
	AliasUint16        = uint16
	AliasUint32        = uint32
	AliasUint64        = uint64
	AliasUintptr       = uintptr
	AliasFloat32       = float32
	AliasFloat64       = float64
	AliasComplex64     = complex64
	AliasComplex128    = complex128
	AliasString        = string
	AliasUnsafePointer = unsafe.Pointer
	AliasByte          = byte
	AliasRune          = rune
)

type BasicAlias struct {
	Bool          AliasBool
	Int           AliasInt
	Int8          AliasInt8
	Int16         AliasInt16
	Int32         AliasInt32
	Int64         AliasInt64
	Uint          AliasUint
	Uint8         AliasUint8
	Uint16        AliasUint16
	Uint32        AliasUint32
	Uint64        AliasUint64
	Uintptr       AliasUintptr
	Float32       AliasFloat32
	Float64       AliasFloat64
	Complex64     AliasComplex64
	Complex128    AliasComplex128
	String        AliasString
	UnsafePointer AliasUnsafePointer
	Byte          AliasByte
	Rune          AliasRune
}

type BasicAliasPtr struct {
	Bool          *AliasBool
	Int           *AliasInt
	Int8          *AliasInt8
	Int16         *AliasInt16
	Int32         *AliasInt32
	Int64         *AliasInt64
	Uint          *AliasUint
	Uint8         *AliasUint8
	Uint16        *AliasUint16
	Uint32        *AliasUint32
	Uint64        *AliasUint64
	Uintptr       *AliasUintptr
	Float32       *AliasFloat32
	Float64       *AliasFloat64
	Complex64     *AliasComplex64
	Complex128    *AliasComplex128
	String        *AliasString
	UnsafePointer *AliasUnsafePointer
	Byte          *AliasByte
	Rune          *AliasRune
}

type (
	AliasArray [6]string
	AliasSlice []string
	AliasMap   map[string]string
)

type ArrayAli struct {
	Name AliasArray
}

type SliceAli struct {
	Name AliasSlice
}

type MapAli struct {
	Name AliasMap
}

type ArrayAliPtr struct {
	Name *AliasArray
}

type SliceAliPtr struct {
	Name *AliasSlice
}

type MapAliPtr struct {
	Name *AliasMap
}

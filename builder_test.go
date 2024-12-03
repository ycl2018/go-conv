package main

import (
	. "go/types"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	*dryRun = false
	*verbose = true
	err := generate("./testdata")
	if err != nil {
		t.Fatal(err)
	}
}

func TestConvertable(t *testing.T) {
	for _, test := range []struct {
		v, t Type
		want bool
	}{
		//{Typ[Int], Typ[Int], true},
		//{Typ[Int], Typ[Float32], true},
		//{Typ[Int], Typ[String], true},
		//{Typ[UntypedInt], Typ[Int], true},
		//{NewSlice(Typ[Int]), NewArray(Typ[Int], 10), true},
		//{NewSlice(Typ[Int]), NewArray(Typ[Uint], 10), false},
		//{NewSlice(Typ[Int]), NewPointer(NewArray(Typ[Int], 10)), true},
		{NewPointer(NewArray(Typ[Int], 10)), NewPointer(NewArray(Typ[Int], 10)), true},
		//{NewSlice(Typ[Int]), NewPointer(NewArray(Typ[Uint], 10)), false},
		//// Untyped string values are not permitted by the spec, so the behavior below is undefined.
		//{Typ[UntypedString], Typ[String], true},
	} {
		if got := ConvertibleTo(test.v, test.t); got != test.want {
			t.Errorf("ConvertibleTo(%v, %v) = %t, want %t", test.v, test.t, got, test.want)
		}
	}
}

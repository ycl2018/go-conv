package a

import "github.com/ycl2018/go-conv/testdata/b"

type NestedSlice struct {
	Slice [][][]*Foo
	Map   map[string]map[string]map[int]*Foo
}

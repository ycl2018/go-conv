package b

type NestedSlice struct {
	Slice [][][]*Foo
	Map   map[string]map[string]map[int]*Foo
}

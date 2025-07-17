package b

type NestedSlice struct {
	Slice      [][][]*Foo
	Map        map[string]map[string]map[int]*Foo
	StringInt  int
	StringInt2 int16
	StringInt3 int32
}

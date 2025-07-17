package a

type NestedSlice struct {
	Slice      [][][]*Foo
	Map        map[string]map[string]map[int]*Foo
	StringInt  string
	StringInt2 string
	StringInt3 string
}

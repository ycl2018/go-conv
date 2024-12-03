package domain

type Category struct {
	CategoryID uint64
	Name       string
	Foo        *Foo
}
type Foo struct {
	Bar *string
}

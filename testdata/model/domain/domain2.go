package domain

type Category struct {
	CategoryID uint64 `storage:"categoryId"`
	Name       string `storage:"name"`
	Foo        *Foo
}
type Foo struct {
	Bar *string
}

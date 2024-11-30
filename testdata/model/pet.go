package model

type Pet struct {
	ID uint `storage:"id"`
	//Category  Category `storage:"category"`
	Name string `storage:"name"`
	//PhotoUrls []string `storage:"photoUrls"`
	Status       *string      `storage:"status"`
	Array        [3]*Category `storage:"children"`
	Slices       []*Category  `storage:"children"`
	Maps         map[string]*Category
	Next         *Pet
	PtrToStruct  *Category
	StructToPtr  Category
	SlicesStruct []Category   `storage:"children"`
	SlicesPtr    []*Category  `storage:"children"`
	ArrayToSlice []*Category  `storage:"children"`
	SliceToArray [3]*Category `storage:"children"`
}

type Category struct {
	CategoryID uint64 `storage:"categoryId"`
	Name       string `storage:"name"`
	Foo        *Foo
}

type Foo struct {
	Bar *string
}

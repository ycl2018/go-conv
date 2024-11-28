package model

type Pet struct {
	ID uint `storage:"id"`
	//Category  Category `storage:"category"`
	Name string `storage:"name"`
	//PhotoUrls []string `storage:"photoUrls"`
	Status    *string      `storage:"status"`
	Children  [3]*Category `storage:"children"`
	Childrens []*Category  `storage:"children"`
}

type Category struct {
	CategoryID uint64 `storage:"categoryId"`
	Name       string `storage:"name"`
	Foo        *Foo
}

type Foo struct {
	Bar *string
}

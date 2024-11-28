package domain

type Pet struct {
	ID uint
	//Category  Category
	Name string
	//PhotoUrls []URL
	Status    *PetStatus
	Children  [3]*Category
	Childrens []*Category `storage:"children"`
}

type URL string

func NewURL(s string) URL {
	return URL(s)
}

func (u URL) String() string {
	return string(u)
}

type Category struct {
	CategoryID uint64 `storage:"categoryId"`
	Name       string `storage:"name"`
	Foo        *Foo
}

type Foo struct {
	Bar *string
}

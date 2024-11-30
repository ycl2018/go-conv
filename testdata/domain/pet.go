package domain

type Pet struct {
	ID uint
	//Category  Category
	Name string
	//PhotoUrls []URL
	Status             *PetStatus
	Array              [3]*Category `storage:"children"`
	Slices             []*Category  `storage:"children"`
	Maps               map[string]*Category
	Next               *Pet
	PtrToStruct        Category
	StructToPtr        *Category
	SlicesStruct       []*Category  `storage:"children"`
	SlicesPtr          []Category   `storage:"children"`
	ArrayToSlice       [3]*Category `storage:"children"`
	SliceToArray       []*Category  `storage:"children"`
	StringConvert2     MyString
	StringConvert      string
	NumberCast         uint64
	ByteSliceToString2 []byte
	ByteSliceToString  string
	MapStringString    MapStringString
}

type MapStringString map[string]string

type MyString string

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

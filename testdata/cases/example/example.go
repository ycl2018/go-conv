package example

// go-conv:generate
// go-conv:copy
var Foo2Bar func(src *Foo) *Bar

type Foo struct {
	Str     string
	Slice   []string
	Map     map[string]string
	Pointer string
	Alias   string
}

type Bar struct {
	Str     string
	Slice   []string
	Map     map[string]string
	Pointer *string
	Alias   StringAlias
}

type StringAlias string

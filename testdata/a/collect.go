package a

type Array struct {
	Name [6]string
}

type Slice struct {
	Name []string
}

type Map struct {
	Name map[string]string
}

type ArrayPtr struct {
	Name *[6]string
}

type SlicePtr struct {
	Name *[]string
}

type MapPtr struct {
	Name *map[string]string
}

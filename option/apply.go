package option

// Option convert option
type Option struct{}

// WithIgnoreFields specify fields of struct to ignore in convert
func WithIgnoreFields(structObj any, fields []string) Option {
	return Option{}
}

// WithIgnoreTypes specify types to ignore in convert
func WithIgnoreTypes(types ...any) Option {
	return Option{}
}

// WithIgnoreIndexes specify Indexes of Slice to ignore in convert
func WithIgnoreIndexes(sliceObj any, indexes []int) Option {
	return Option{}
}

// WithIgnoreKeys specify Keys of Map to ignore in convert
func WithIgnoreKeys(mapObj any, keys []any) Option {
	return Option{}
}

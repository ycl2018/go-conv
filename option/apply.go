package option

// Option convert option
type Option struct{}

// WithIgnoreFields specify fields of struct to ignore in convert
func WithIgnoreFields(structObj any, fields []string, path ...string) Option {
	return Option{}
}

// WithIgnoreTypes specify types to ignore in convert
func WithIgnoreTypes(types []any, path ...string) Option {
	return Option{}
}

// WithTransformer specify transformer function on type T. if path are set, transformer only apply to specified path.
// Or else all types that matched with transformer will be applied.
func WithTransformer[T, V any](transformer func(T) V, path ...string) Option {
	return Option{}
}

// WithFilter specify filter function on type T. if path are set, filter only apply to specified path.
// Or else all types that matched with filter will be applied.
func WithFilter[T any](filter func(T) T, path ...string) Option {
	return Option{}
}

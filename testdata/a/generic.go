package a

type Generic[K comparable, V any] struct {
	Map   map[K]V
	Slice []V
	Array [3]V
	K     K
	V     V
}

package typex

func ListOf[T any](v T, vv ...T) []T {
	return append(append(make([]T, 0, len(vv)+1), v), vv...)
}

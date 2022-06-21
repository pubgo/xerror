package typex

func Zero[T any]() (ret T) {
	return
}

func Nil[T any]() (ret *T) {
	return
}

func Ptr[T any](a T) *T {
	return &a
}

func ListOf[T any](v T, vv ...T) []T {
	return append(append(make([]T, 0, len(vv)+1), v), vv...)
}

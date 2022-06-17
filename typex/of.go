package typex

func Zero[T any]() (ret T) {
	return
}

func Of[T any](v T, vv ...T) List[T] {
	return append(append(make([]T, 0, len(vv)+1), v), vv...)
}

type List[T any] []T

func (a *List[T]) Append(data ...T) List[T] {
	return append(append(make([]T, 0, len(*a)+len(data)), *a...), data...)
}

func (a *List[T]) Range(fn func(v T)) {
	for _, v := range *a {
		fn(v)
	}
}

func (a *List[T]) Map(fn func(v T) T) {
	for i, v := range *a {
		(*a)[i] = fn(v)
	}
}

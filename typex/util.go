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

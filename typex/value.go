package typex

type Result[T any] struct {
	err error
	val T
}

func (v Result[T]) IsErr() bool { return v.err != nil }
func (v Result[T]) Get() T      { return v.val }
func (v Result[T]) Err() error  { return v.err }

func OK[T any](val T, err error) Result[T] {
	return Result[T]{val: val, err: err}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

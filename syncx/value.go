package syncx

type Value[T any] struct {
	err error
	val T
}

func (v Value[T]) IsErr() bool { return v.err != nil }
func (v Value[T]) Val() T      { return v.val }
func (v Value[T]) Err() error  { return v.err }

func OK[T any](val T, err ...error) Value[T] {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return Value[T]{val: val, err: e}
}

func Err[T any](err error) Value[T] {
	return Value[T]{err: err}
}

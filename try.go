package funk

import "github.com/pubgo/funk/internal/utils"

func checkFn(fn interface{}) {
	if fn == nil {
		panic("[fn] should not be nil")
	}
}

func Try(fn func()) (gErr error) {
	checkFn(fn)

	defer Recovery(func(err XErr) {
		gErr = err.WrapF("fn=>%s", utils.CallerWithFunc(fn))
	})

	fn()
	return
}

func TryWith(gErr *error, fn func()) {
	checkFn(fn)

	defer Recovery(func(err XErr) {
		*gErr = err.WrapF("fn=>%s", utils.CallerWithFunc(fn))
	})

	fn()
}

func TryCatch(fn func(), catch func(err error)) {
	checkFn(fn)
	checkFn(catch)

	defer Recovery(func(err XErr) {
		catch(err.WrapF("fn=>%s", utils.CallerWithFunc(fn)))
	})

	fn()
}

func TryThrow(fn func()) {
	checkFn(fn)

	defer RecoverAndRaise(func(err XErr) XErr {
		return err.WrapF("fn=>", utils.CallerWithFunc(fn))
	})

	fn()
}

func TryVal[T any](fn func() (T, error), cache func(err error)) T {
	defer Recovery(func(err XErr) {
		cache(err.WrapF("fn=>", utils.CallerWithFunc(fn)))
	})

	checkFn(fn)

	val, err := fn()
	if err == nil {
		return val
	}
	cache(err)
	return val
}

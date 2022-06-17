package xerror

import (
	"github.com/pubgo/xerror/internal/utils"
)

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

func TryErr(gErr *error, fn func()) {
	checkFn(fn)

	defer Recovery(func(err XErr) {
		*gErr = err.WrapF("fn=>%s", utils.CallerWithFunc(fn))
	})

	fn()

	return
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

func TryVal[T any](fn func() (T, error), cache func(val T)) (gErr error) {
	defer RecoverErr(&gErr, func(err XErr) XErr {
		return err.WrapF("fn=>", utils.CallerWithFunc(fn))
	})

	checkFn(fn)

	val, err := fn()
	if err != nil {
		return err
	}

	cache(val)
	return nil
}

package funk

import (
	"github.com/pubgo/funk/internal/utils"
	"k8s.io/klog/v2"
)

func Try(fn func()) (gErr error) {
	Assert(fn == nil, "[fn] is nil")

	defer RecoverErr(&gErr, func(err XErr) XErr {
		return err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})

	fn()
	return
}

func TryWith(err *error, fn func()) {
	Assert(fn == nil, "[fn] is nil")

	defer RecoverErr(err, func(err XErr) XErr {
		return err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})

	fn()
}

func TryAndLog(fn func(), catch ...func(err XErr) XErr) {
	Assert(fn == nil, "[fn] is nil")

	defer Recovery(func(err XErr) {
		if len(catch) > 0 {
			err = catch[0](err)
		}

		err = err.WrapF("fn=%s", utils.CallerWithFunc(fn))
		err.DebugPrint()
		klog.Error(err.Error(), " ", err)
	})

	fn()
}

func TryCatch(fn func(), catch func(err error)) {
	Assert(fn == nil, "[fn] is nil")
	Assert(catch == nil, "[catch] is nil")

	defer Recovery(func(err XErr) {
		catch(err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
	})

	fn()
}

func TryThrow(fn func()) {
	Assert(fn == nil, "[fn] is nil")

	defer RecoverAndRaise(func(err XErr) XErr {
		return err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})

	fn()
}

func TryRet[T any](fn func() (T, error), cache func(err error)) T {
	Assert(fn == nil, "[fn] is nil")

	defer Recovery(func(err XErr) {
		cache(err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
	})

	val, err := fn()
	if err == nil {
		return val
	}
	cache(err)
	return val
}

package funk

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/logx"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/xerr"
)

func Try(fn func() error) (gErr error) {
	assert.If(fn == nil, "[fn] is nil")

	defer recovery.Err(&gErr, func(err xerr.XErr) xerr.XErr {
		return err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})

	return fn()
}

func TryWith(err *error, fn func() error) {
	assert.If(fn == nil, "[fn] is nil")

	defer recovery.Err(err, func(err xerr.XErr) xerr.XErr {
		return err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})

	*err = fn()
}

func TryAndLog(fn func(), catch ...func(err xerr.XErr) xerr.XErr) {
	assert.If(fn == nil, "[fn] is nil")

	defer recovery.Recovery(func(err xerr.XErr) {
		if len(catch) > 0 {
			err = catch[0](err)
		}

		err = err.WrapF("fn=%s", utils.CallerWithFunc(fn))
		logx.Error(err, "log panic func")
	})

	fn()
}

func TryCatch(fn func() error, catch func(err xerr.XErr)) {
	assert.If(fn == nil, "[fn] is nil")
	assert.If(catch == nil, "[catch] is nil")

	defer recovery.Recovery(func(err xerr.XErr) {
		catch(err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
	})

	catch(xerr.WrapXErr(fn()))
}

func TryThrow(fn func()) {
	assert.If(fn == nil, "[fn] is nil")

	defer recovery.Raise(func(err xerr.XErr) xerr.XErr {
		return err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})

	fn()
}

func TryRet[T any](fn func() (T, error), cache func(err xerr.XErr)) T {
	assert.If(fn == nil, "[fn] is nil")

	defer recovery.Recovery(func(err xerr.XErr) {
		cache(err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
	})

	val, err := fn()
	if err == nil {
		return val
	}
	cache(xerr.WrapXErr(err))
	return val
}

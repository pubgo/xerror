package recovery

import (
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/xerr"
)

func CatchErr(gErr *error, fns ...func(err xerr.XErr) xerr.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	xerr.ParseErr(&err, val)
	if err == nil {
		return
	}

	err1 := xerr.WrapXErr(err)
	if len(fns) > 0 {
		*gErr = fns[0](err1)
		return
	}
	*gErr = err1
}

func Err(gErr *error, fns ...func(err xerr.XErr) xerr.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	xerr.ParseErr(&err, val)
	if err == nil {
		return
	}

	err1 := xerr.WrapXErr(err)
	if len(fns) > 0 {
		*gErr = fns[0](err1)
		return
	}
	*gErr = err1
}

func Raise(fns ...func(err xerr.XErr) xerr.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	xerr.ParseErr(&err, val)
	if err == nil {
		return
	}

	err1 := xerr.WrapXErr(err)
	if len(fns) > 0 {
		panic(fns[0](err1))
	}
	panic(err1)
}

func Recovery(fn func(err xerr.XErr)) {
	assert.If(fn == nil, "[fn] should not be nil")

	val := recover()
	if val == nil {
		return
	}

	var err error
	xerr.ParseErr(&err, val)
	if err == nil {
		return
	}

	fn(xerr.WrapXErr(err))
}

func Exit() {
	val := recover()
	if val == nil {
		return
	}

	var err error
	xerr.ParseErr(&err, val)
	if err == nil {
		return
	}

	xerr.WrapXErr(err).DebugPrint()
	os.Exit(1)
}

func Log() {
	val := recover()
	if val == nil {
		return
	}

	var err error
	xerr.ParseErr(&err, val)
	if err == nil {
		return
	}

	xerr.WrapXErr(err).DebugPrint()
	debug.PrintStack()
}

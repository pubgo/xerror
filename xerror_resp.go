package xerror

import (
	"os"

	"github.com/pubgo/xerror/xerror_envs"
	"github.com/pubgo/xerror/xerror_util"
)

func RespErr(err *error) {
	handleErr(err, recover())
	if isErrNil(*err) {
		*err = nil
	}
}

func RespJson() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, xerrorOptions{}).Stack(true))
}

func RespDebug() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, xerrorOptions{}).p())
	xerror_envs.PrintStackVal()
}

func RespRaise(format string, a ...interface{}) {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	With(WithCaller(5)).PanicF(err, format, a...)
}

// Resp
func Resp(f func(err XErr)) {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	if err, ok := err.(XErr); ok {
		f(err.(XErr))
		return
	}
	f(&xerror{Cause1: err, Caller: xerror_util.CallerWithDepth(xerror_envs.CallDepthVal() + 1)})
}

func RespExit() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, xerrorOptions{}).p())
	xerror_envs.PrintStackVal()
	os.Exit(1)
}

func RespGoroutine(name ...string) {
	nm := "__xerror"
	if len(name) > 0 {
		nm = name[0]
	}

	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	goroutineErrs <- &goroutineErrEvent{name: nm, err: handle(err, xerrorOptions{})}
}

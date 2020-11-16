package xerror

import (
	"os"

	"github.com/pubgo/xerror/internal/wrapper"
	"github.com/pubgo/xerror/xerror_util"
)

func RespErr(err *error) {
	handleErr(err, recover())
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
	wrapper.PrintStack()
}

// Resp
func Resp(f func(err XErr)) {
	var err error
	handleErr(&err, recover())
	if err == nil {
		return
	}

	if err, ok := err.(XErr); ok {
		f(err.(XErr))
		return
	}
	f(&xerror{Cause1: err, Caller: xerror_util.CallerWithDepth(wrapper.CallDepth() + 1)})
}

func RespExit() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, xerrorOptions{}).p())
	wrapper.PrintStack()
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

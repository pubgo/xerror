package xerror

import (
	"os"

	"github.com/pubgo/xerror/xerror_envs"
	"github.com/pubgo/xerror/xerror_util"
)

func RespErr(err *error) {
	handleRecover(err, recover())
	if isErrNil(*err) {
		*err = nil
	}
}

func RespJson() {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, xerrorOptions{}).Stack(true))
}

func RespDebug() {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, xerrorOptions{}).p())
	xerror_envs.PrintStackVal()
}

func RespRaise(fn func(err XErr) error) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	panic(&xerror{Cause1: fn(err.(XErr)), Caller: xerror_util.CallerWithFunc(fn)})
}

// Resp
func Resp(f func(err XErr)) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	f(&xerror{Cause1: err, Caller: xerror_util.CallerWithFunc(f)})
}

func RespExit() {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, xerrorOptions{}).p())
	xerror_envs.PrintStackVal()
	os.Exit(1)
}

package xerror

import (
	"os"

	"github.com/pubgo/xerror/xerror_abc"
	"github.com/pubgo/xerror/xerror_util"
)

func RespErr(err *error) {
	handleRecover(err, recover())
	if isErrNil(*err) {
		*err = nil
	}
}

func RespDebug() {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, options{}).p())
	printStack()
}

func RespRaise(fn func(err xerror_abc.XErr) error) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	panic(fn(&xerror{Cause1: err, Caller: xerror_util.CallerWithFunc(fn)}))
}

// Resp
func Resp(f func(err xerror_abc.XErr)) {
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

	p(handle(err, options{}).p())
	printStack()
	os.Exit(1)
}

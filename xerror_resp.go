package xerror

import (
	"fmt"
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

func RespDebug(args ...interface{}) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, options{msg: fmt.Sprint(args...)}).p())
	printStack()
}

func RespRaise(fn func(err xerror_abc.XErr) error) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	err1 := &xerror{Cause1: err, Caller: xerror_util.CallerWithFunc(fn)}
	if fn == nil {
		panic(err1)
	}
	panic(fn(err1))
}

// Resp
func Resp(fn func(err xerror_abc.XErr)) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	Assert(fn == nil, "[fn] should not be nil")
	fn(&xerror{Cause1: err, Caller: xerror_util.CallerWithFunc(fn)})
}

func RespExit(args ...interface{}) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, options{msg: fmt.Sprint(args...)}).p())
	printStack()
	os.Exit(1)
}

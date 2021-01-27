package xerror

import (
	"fmt"
	"os"

	"github.com/pubgo/xerror/xerror_util"
)

func RespErr(err *error) { handleRecover(err, recover()) }

func RespDebug(args ...interface{}) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }).stackString())
	printStack()
}

func Raise(fns ...func(err XErr) error) { RespRaise(fns...) }
func RespRaise(fns ...func(err XErr) error) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	err1 := &xerror{Cause1: err}
	if len(fns) > 0 {
		err1.Caller = [2]string{xerror_util.CallerWithFunc(fns[0])}
		panic(fns[0](err1))
	}

	panic(err1)
}

// Resp
func Resp(fn func(err XErr)) {
	Assert(fn == nil, "[fn] should not be nil")

	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	fn(&xerror{Cause1: err, Caller: [2]string{xerror_util.CallerWithFunc(fn)}})
}

func RespExit(args ...interface{}) {
	var err error
	handleRecover(&err, recover())
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }).stackString())
	printStack()
	os.Exit(1)
}

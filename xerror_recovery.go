package xerror

import (
	"fmt"
	"os"

	"github.com/pubgo/xerror/internal/utils"
)

func RespErr(err *error) {
	val := recover()
	if val == nil {
		return
	}

	handleRecover(err, val)
}

func RespDebug(args ...interface{}) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }).stackString())
	printStack()
}

func Raise(fns ...func(err XErr) error) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	err1 := &xerror{Err: err}
	if len(fns) > 0 {
		err1.Caller = [2]string{utils.CallerWithFunc(fns[0])}
		panic(fns[0](err1))
	}

	panic(err1)
}

func RespRaise(fns ...func(err XErr) error) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	err1 := &xerror{Err: err}
	if len(fns) > 0 {
		err1.Caller = [2]string{utils.CallerWithFunc(fns[0])}
		panic(fns[0](err1))
	}

	panic(err1)
}

// Resp
func Resp(fn func(err XErr)) {
	Assert(fn == nil, "[fn] should not be nil")

	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	fn(&xerror{Err: err, Caller: [2]string{utils.CallerWithFunc(fn)}})
}

func RespExit(args ...interface{}) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }).stackString())
	printStack()
	os.Exit(1)
}
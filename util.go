package xerror

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pubgo/xerror/xerror_core"
	"github.com/pubgo/xerror/xerror_util"
)

func handleRecover(err *error, err1 interface{}) {
	if err1 == nil {
		return
	}

	switch err1 := err1.(type) {
	case error:
		*err = err1
	case string:
		*err = errors.New(err1)
	default:
		*err = WrapF(ErrType, fmt.Sprintf("%#v", err1))
	}
}

func handle(err error, opts options) *xerror {
	err2 := &xerror{}
	err2.Msg = opts.msg
	err2.Caller = xerror_util.CallerWithDepth(xerror_core.Conf.CallDepth + 1 + opts.depth)
	switch err := err.(type) {
	case *xerrorBase, *xerror, *combine, error:
		err2.Cause1 = err
	default:
		err2.Cause1 = WrapF(ErrType, fmt.Sprintf("%#v", err))
	}

	return err2
}

func isErrNil(err error) bool {
	return err == nil || err == ErrDone || Unwrap(err) == ErrDone
}

func trans(err error) []*xerror {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *xerrorBase:
		return []*xerror{{
			Msg:    err.Msg,
			Caller: err.Caller,
		}}
	case *xerror:
		return []*xerror{err}
	case *combine:
		return *err
	default:
		return nil
	}
}

func Unwrap(err error) error {
	for {
		u, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			return err
		}
		err = u.Unwrap()
	}
}

func p(a ...interface{}) { _, _ = fmt.Fprintln(os.Stderr, a...) }
func printStack() {
	if xerror_core.Conf.PrintStack {
		debug.PrintStack()
	}
}

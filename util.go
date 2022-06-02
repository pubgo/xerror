package xerror

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pubgo/xerror/internal/utils"
	"github.com/pubgo/xerror/xerror_core"
)

func isErrNil(err error) bool { return err == nil }
func p(a ...interface{})      { _, _ = fmt.Fprintln(os.Stderr, a...) }

func handleRecover(err *error, val interface{}) {
	if val == nil {
		return
	}

	// 自定义error检测
	for _, h := range xerror_core.Handlers() {
		if _err := h(val); _err != nil {
			*err = _err
			return
		}
	}

	switch _val := val.(type) {
	case error:
		*err = _val
	case string:
		*err = errors.New(_val)
	default:
		*err = fmt.Errorf("%#v\n", _val)
	}
}

func handle(err error, fns ...func(err *xerror)) *xerror {
	err1 := &xerror{Err: err, Msg: err.Error()}
	err1.Caller = []string{
		utils.CallerWithDepth(xerror_core.Conf.CallDepth + 2),
		utils.CallerWithDepth(xerror_core.Conf.CallDepth + 3),
	}

	if len(fns) > 0 {
		fns[0](err1)
	}

	return err1
}

func trans(err error) *xerror {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *xerror:
		return err
	case interface{ Unwrap() error }:
		if err.Unwrap() == nil {
			return nil
		}
		return &xerror{Err: err.Unwrap(), Msg: err.Unwrap().Error()}
	default:
		return nil
	}
}

func printStack() {
	if !xerror_core.Conf.PrintStack {
		return
	}

	debug.PrintStack()
}

func Cause(err error) error {
	if isErrNil(err) {
		return nil
	}

	for {
		rErr := errors.Unwrap(err)
		if isErrNil(rErr) {
			return err
		}

		err = rErr
	}
}

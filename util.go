package xerror

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pubgo/xerror/xerror_envs"
	"github.com/pubgo/xerror/xerror_util"
)

func handleErr(err *error, _err interface{}) {
	if _err == nil {
		return
	}

	switch _err := _err.(type) {
	case error:
		*err = _err
	case string:
		*err = errors.New(_err)
	default:
		*err = WrapF(ErrUnknownType, fmt.Sprintf("%#v", _err))
	}
}

func handle(err error, opts xerrorOptions) *xerror {
	err2 := &xerror{}
	err2.Msg = opts.msg
	err2.Caller = xerror_util.CallerWithDepth(xerror_envs.CallDepthVal() + 1 + opts.depth)
	switch e := err.(type) {
	case *xerrorBase:
		err2.Cause1 = e
	case *xerror:
		err2.Cause1 = e
	case *xerrorCombine:
		err2.Cause1 = e
	case error:
		err2.Cause1 = &xerrorBase{Code: unwrap(e).Error(), Msg: fmt.Sprintf("%+v", e)}
	default:
		err2.Cause1 = &xerrorBase{Code: ErrUnknownType.Error(), Msg: fmt.Sprintf("%+v", e)}
	}

	return err2
}

func isErrNil(err error) bool {
	return err == nil || err == ErrDone
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
	case *xerrorCombine:
		return *err
	default:
		return nil
	}
}

func unwrap(err error) error {
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

func p(a ...interface{}) {
	_, _ = os.Stderr.WriteString(fmt.Sprintln(a...))
}

func PrintStack() {
	if xerror_envs.PrintStackVal() {
		debug.PrintStack()
	}
}

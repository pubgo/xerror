package xerror

import (
	"errors"
	"fmt"
	"github.com/pubgo/xerror/xerror_envs"
	"github.com/pubgo/xerror/xerror_util"
	"os"
)

func handleRecover(err *error, err1 interface{}) {
	if err1 == nil {
		return
	}

	switch _err := err1.(type) {
	case error:
		*err = _err
	case string:
		*err = errors.New(_err)
	default:
		*err = WrapF(ErrType, fmt.Sprintf("%#v", _err))
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
	case *combine:
		err2.Cause1 = e
	case error:
		err2.Cause1 = &xerrorBase{Code: unwrap(e).Error(), Msg: fmt.Sprintf("%+v", e)}
	default:
		err2.Cause1 = &xerrorBase{Code: ErrType.Error(), Msg: fmt.Sprintf("%+v", e)}
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
	case *combine:
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

func p(a ...interface{}) { _, _ = fmt.Fprintln(os.Stderr, a...) }
func printStack()        { xerror_util.PrintDebug() }

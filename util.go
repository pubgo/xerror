package xerror

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pubgo/xerror/internal/utils"
	"github.com/pubgo/xerror/xerror_core"
)

func handleRecover(err *error, val interface{}) {
	switch val := val.(type) {
	case error:
		*err = val
	case string:
		*err = errors.New(val)
	default:
		*err = fmt.Errorf("%#v\n", val)
	}
}

func handle(err error, fns ...func(err *xerror)) *xerror {
	err2 := &xerror{}
	err2.Caller[0] = utils.CallerWithDepth(xerror_core.Conf.CallDepth + 2)
	err2.Caller[1] = utils.CallerWithDepth(xerror_core.Conf.CallDepth + 3)
	switch err := err.(type) {
	case *xerrorBase, *xerror, *multiError, error:
		err2.Err = err
	default:
		err2.Err = WrapF(ErrType, fmt.Sprintf("%#v", err))
	}

	if len(fns) > 0 {
		fns[0](err2)
	}

	return err2
}

func isErrNil(err error) bool {
	return err == nil
}

func trans(err error) []*xerror {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *xerrorBase:
		return []*xerror{{
			Msg:    err.Msg,
			Caller: [2]string{err.Caller},
		}}
	case *xerror:
		return []*xerror{err}
	case *multiError:
		var errs []*xerror
		for i := range err.errors {
			errs = append(errs, &xerror{Err: err.errors[i]})
		}
		return errs
	default:
		return nil
	}
}

func p(a ...interface{}) { _, _ = fmt.Fprintln(os.Stderr, a...) }
func printStack() {
	if !xerror_core.Conf.PrintStack {
		return
	}

	debug.PrintStack()
}

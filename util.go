package funk

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/funkonf"
	"github.com/pubgo/funk/internal/utils"
)

func isErrNil(err error) bool { return err == nil }
func p(a ...interface{})      { _, _ = fmt.Fprintln(os.Stderr, a...) }

func handleRecover(err *error, val interface{}) {
	if val == nil {
		return
	}

	switch _val := val.(type) {
	case error:
		*err = _val
	case string:
		*err = errors.New(_val)
	default:
		*err = fmt.Errorf("%#v", _val)
	}

	*err = handle(*err)
}

func handle(err error, fns ...func(err *xerror)) *xerror {
	err1 := &xerror{Err: err}
	if _, ok := err.(XErr); !ok {
		for i := 0; ; i++ {
			var cc = utils.CallerWithDepth(callStackDepth + i)
			if cc == "" {
				break
			}
			err1.Caller = append(err1.Caller, cc)
		}
	} else {
		err1.Caller = []string{
			utils.CallerWithDepth(callStackDepth + 2),
		}
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
			return &xerror{Detail: fmt.Sprintf("%#v", err)}
		}
		return &xerror{Err: err.Unwrap(), Msg: err.Unwrap().Error()}
	default:
		return &xerror{Msg: err.Error(), Detail: fmt.Sprintf("%#v", err)}
	}
}

func printStack() {
	if !funkonf.Conf.PrintStack {
		return
	}

	debug.PrintStack()
}

func Last[T any](args []T) (t T) {
	if len(args) == 0 {
		return
	}

	return args[len(args)-1]
}

func Ternary[T any](ok bool, a T, b T) T {
	if ok {
		return a
	}
	return b
}

func If(ok bool, fn func()) {
	if ok {
		fn()
	}
}

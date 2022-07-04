package xerr

import (
	"errors"
	"fmt"
	"os"

	"github.com/pubgo/funk/internal/utils"
)

func isErrNil(err error) bool { return err == nil }
func p(a ...interface{})      { _, _ = fmt.Fprintln(os.Stderr, a...) }

func ParseErr(err *error, val interface{}) {
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

	*err = WrapXErr(*err)
}

func WrapXErr(err error, fns ...func(err *Xerror)) *Xerror {
	err1 := &Xerror{Err: err}
	if _, ok := err.(XErr); !ok {
		for i := 0; ; i++ {
			var cc = utils.CallerWithDepth(CallStackDepth + i)
			if cc == "" {
				break
			}
			err1.Caller = append(err1.Caller, cc)
		}
	} else {
		err1.Caller = []string{
			utils.CallerWithDepth(CallStackDepth + 2),
		}
	}

	if len(fns) > 0 {
		fns[0](err1)
	}

	return err1
}

func trans(err error) *Xerror {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *Xerror:
		return err
	case interface{ Unwrap() error }:
		if err.Unwrap() == nil {
			return &Xerror{Detail: fmt.Sprintf("%#v", err)}
		}
		return &Xerror{Err: err.Unwrap(), Msg: err.Unwrap().Error()}
	default:
		return &Xerror{Msg: err.Error(), Detail: fmt.Sprintf("%#v", err)}
	}
}

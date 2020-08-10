package xerror

import (
	"errors"
	"fmt"
	"github.com/pubgo/xerror/internal/wrapper"
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
		*err = WrapF(ErrUnknownType, fmt.Sprintf("%+v", _err))
	}
}

func handle(err error, msg string, args ...interface{}) *xerror {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	err2 := &xerror{}
	err2.Msg = msg
	err2.Caller = xerror_util.CallerWithDepth(wrapper.CallDepth() + 1)
	switch e := err.(type) {
	case *xerrorBase:
		err2.Cause1 = e
	case *xerror:
		err2.Cause1 = e
	case error:
		err2.Cause1 = New(unwrap(e).Error(), fmt.Sprintf("%+v", e))
	default:
		err2.Cause1 = New(ErrUnknownType.Error(), fmt.Sprintf("%+v", e))
	}

	return err2
}

func isErrNil(err error) bool {
	return err == nil || err == ErrDone
}

func trans(err error) *xerror {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *xerrorBase:
		return &xerror{
			Msg:    err.Msg,
			Caller: err.Caller,
		}
	case *xerror:
		return err
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

package xerror

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/pubgo/xerror/internal/utils"
	"github.com/pubgo/xerror/xerror_core"
)

var xerrorTyp = reflect.TypeOf(&xerror{})
var xerrorBaseTyp = reflect.TypeOf(&xerrorBase{})

func isErrNil(err error) bool { return err == nil }
func p(a ...interface{})      { _, _ = fmt.Fprintln(os.Stderr, a...) }

// Parse parse error to xerror
func Parse(err error) XErr {
	if isErrNil(err) {
		return nil
	}

	return handle(err)
}

// ParseWith parse error to xerror
func ParseWith(err error, fn func(err XErr)) {
	if isErrNil(err) {
		return
	}

	fn(handle(err))
}

func IsXErr(err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case *xerrorBase:
		return true
	case *xerror:
		return true
	case *multiError:
		return true
	default:
		return false
	}
}

func handleRecover(err *error, val interface{}) {
	if val == nil {
		return
	}

	switch val1 := val.(type) {
	case error:
		*err = val1
	case string:
		*err = errors.New(val1)
	default:
		*err = fmt.Errorf("%#v\n", val1)
	}
}

func handle(err error, fns ...func(err *xerror)) *xerror {
	err1 := &xerror{}
	err1.Caller[0] = utils.CallerWithDepth(xerror_core.Conf.CallDepth + 2)
	err1.Caller[1] = utils.CallerWithDepth(xerror_core.Conf.CallDepth + 3)
	switch err := err.(type) {
	case *xerrorBase, *xerror, *multiError, error:
		err1.Err = err
	default:
		err1.Err = WrapF(ErrType, fmt.Sprintf("%#v", err))
	}

	if len(fns) > 0 {
		fns[0](err1)
	}

	return err1
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

func printStack() {
	if !xerror_core.Conf.PrintStack {
		return
	}

	debug.PrintStack()
}

func As(err error, target interface{}) bool {
	if target == nil || err == nil {
		return false
	}

	val := reflect.ValueOf(target)
	typ := val.Type()

	// target must be a non-nil pointer
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}

	// *target must be interface or implement error
	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return false
	}

	errType := reflect.TypeOf(err)
	for {
		if errType != xerrorTyp && errType != xerrorBaseTyp && reflect.TypeOf(err).AssignableTo(typ.Elem()) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}

		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}

		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
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

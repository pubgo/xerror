package xerror

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/pubgo/xerror/xerror_abc"
)

// PanicErrs combine multiple errors
func PanicErrs(errs ...error) {
	if len(errs) == 0 {
		return
	}

	panic(Combine(errs...))
}

// Parse parse error to xerror
func Parse(err error) xerror_abc.XErr {
	if isErrNil(err) {
		return nil
	}

	return handle(err)
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
	default:
		return false
	}
}

func Try(fn func()) (err error) {
	Assert(fn == nil, "[fn] should not be nil")

	defer RespErr(&err)
	fn()
	return
}

func Done() { panic(ErrDone) }

func Panic(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }))
}

func PanicF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(handle(err, func(err *xerror) { err.Msg = fmt.Sprintf(msg, args...) }))
}

func Wrap(err error, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	return handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) })
}

func WrapF(err error, msg string, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	return handle(err, func(err *xerror) { err.Msg = fmt.Sprintf(msg, args...) })
}

// PanicErr
func PanicErr(d1 interface{}, err error) interface{} {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err))
}

func PanicBytes(d1 []byte, err error) []byte {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err))
}

func PanicStr(d1 string, err error) string {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err))
}

func Exit(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }).stackString())
	printStack()
	os.Exit(1)
}

// ExitF
func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprintf(msg, args...) }).stackString())
	printStack()
	os.Exit(1)
}

// ExitErr
func ExitErr(dat interface{}, err error) interface{} {
	if isErrNil(err) {
		return dat
	}

	p(handle(err).stackString())
	printStack()
	os.Exit(1)
	return nil
}

func Is(err, target error) bool { return errors.Is(err, target) }
func Unwrap(err error) error    { return errors.Unwrap(err) }

var xerrorTyp = reflect.TypeOf(&xerror{})
var xerrorBaseTyp = reflect.TypeOf(&xerrorBase{})

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

		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func Cause(err error) error {
	for {
		err1 := Unwrap(err)
		if err1 == nil {
			return err
		}
		err = err1
	}
}

package xerror

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
)

type XErr interface {
	error
	As(err interface{}) bool
	Is(err error) bool
	Unwrap() error
	Cause() error
	Code() string
	Format(s fmt.State, verb rune)
	Stack() string
}

func New(code string, ms ...string) *xerrorBase {
	var msg string
	if len(ms) > 0 {
		msg = ms[0]
	}

	xw := &xerrorBase{}
	xw.Code = code
	xw.Msg = msg
	xw.Caller = callerWithDepth(callDepth)

	return xw
}

func Try(fn func()) (err error) {
	defer Resp(func(_err XErr) {
		err = handle(_err, "")
		err.(*xerror).Caller = callerWithFunc(reflect.ValueOf(fn))
	})
	fn()
	return
}

func RespErr(err *error) {
	handleErr(err, recover())
}

// Resp
func Resp(f func(err XErr)) {
	var err error
	handleErr(&err, recover())
	if err != nil {
		f(err.(XErr))
	}
}

func RespExit() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	fmt.Println(handle(err, "").(*xerror).p())
	debug.PrintStack()
	os.Exit(1)
}

func Panic(err error) {
	if isErrNil(err) {
		return
	}

	panic(handle(err, ""))
}

func PanicF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}
	panic(handle(err, msg, args...))
}

func Wrap(err error) error {
	if isErrNil(err) {
		return nil
	}

	return handle(err, "")
}

func WrapF(err error, msg string, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	return handle(err, msg, args...)
}

// PanicErr
func PanicErr(d1 interface{}, err error) interface{} {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err, ""))
}

func PanicBytes(d1 []byte, err error) []byte {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err, ""))
}

func PanicStr(d1 string, err error) string {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err, ""))
}

func PanicFile(d1 *os.File, err error) *os.File {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err, ""))
}

func PanicResponse(d1 *http.Response, err error) *http.Response {
	if isErrNil(err) {
		return d1
	}

	panic(handle(err, ""))
}

// ExitErr
func ExitErr(_ interface{}, err error) {
	if isErrNil(err) {
		return
	}

	fmt.Println(handle(err, "").(*xerror).p())
	debug.PrintStack()
	os.Exit(1)
}

// ExitF
func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	fmt.Println(handle(err, msg, args...).(*xerror).p())
	debug.PrintStack()
	os.Exit(1)
}

func Exit(err error) {
	if isErrNil(err) {
		return
	}

	fmt.Println(handle(err, "").(*xerror).p())
	debug.PrintStack()
	os.Exit(1)
}

func Unwrap(err error) error {
	for !isErrNil(err) {
		wrap, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = wrap.Unwrap()
	}
	return err
}

func Is(err, target error) bool {
	if isErrNil(target) {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		if err = Unwrap(err); isErrNil(err) {
			return false
		}
	}
}

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

func As(err error, target interface{}) bool {
	if target == nil {
		return false
	}

	val := reflect.ValueOf(target)
	typ := val.Type()

	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}

	if e := typ.Elem(); e.Kind() != reflect.Interface && !typ.Implements(errorType) {
		return false
	}

	targetType := typ.Elem()
	for !isErrNil(err) {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

// Cause returns the underlying cause of the error, if possible.
func Cause(err error) error {
	for !isErrNil(err) {
		cause, ok := err.(interface{ Cause() error })
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

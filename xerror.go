package xerror

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
)

type XErr interface {
	New(code, msg string) XErr
	XRErr
}

type XRErr interface {
	error
	fmt.Formatter
	As(err interface{}) bool
	Is(err error) bool
	Unwrap() error
	Code() string
	Detail() string
	Reset()
}

func New(code, msg string) XErr {
	return &xerror{Code1: code, Msg: msg}
}

func Try(fn func() error) (err error) {
	defer Resp(func(_err XRErr) {
		err = handle(_err, "")
		err.(*xerror).Caller = callerWithFunc(reflect.ValueOf(fn))
	})
	err = fn()
	if isErrNil(err) {
		return
	}
	panic(handle(err, ""))
}

func RespErr(err *error) {
	handleErr(err, recover())
}

// Resp
func Resp(f func(err XRErr)) {
	var err error
	handleErr(&err, recover())
	if err != nil {
		f(err.(XErr))
	}
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

	logger.Println(handle(err, "").Error())
	if Debug {
		debug.PrintStack()
	}
	os.Exit(1)
}

// ExitF
func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	logger.Println(handle(err, msg, args...).Error())
	if Debug {
		debug.PrintStack()
	}
	os.Exit(1)
}

func Exit(err error) {
	if isErrNil(err) {
		return
	}

	logger.Println(handle(err, "").Error())
	if Debug {
		debug.PrintStack()
	}
	os.Exit(1)
}

// ext from errors
var (
	UnWrap = errors.Unwrap
	Is     = errors.Is
	As     = func(err error, target interface{}) bool {
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
		for err != nil {
			if reflect.TypeOf(err).AssignableTo(targetType) {
				val.Elem().Set(reflect.ValueOf(err))
				return true
			}
			if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
				return true
			}
			err = UnWrap(err)
		}
		return false
	}
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

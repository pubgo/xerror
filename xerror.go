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
	error
	As(err interface{}) bool
	Is(err error) bool
	Unwrap() error
	Code() string
	Stack() string
	Reset()
}

type XErrBase interface {
	error
	New(code string, ms ...string) XErrBase
	Code() string
}

type xerrorBase struct {
	*xerror
}

func (t *xerrorBase) New(code string, ms ...string) XErrBase {
	var msg string
	if len(ms) == 1 {
		msg = ms[0]
	}

	code = t.Code1 + ": " + code
	xw := &xerrorBase{xerror: new(xerror)}
	xw.Code1 = code
	xw.Msg = msg
	xw.xrr = errors.New(code)
	xw.Caller = callerWithDepth(callDepth)

	return xw
}

func New(code string, ms ...string) XErrBase {

	var msg string
	if len(ms) == 1 {
		msg = ms[0]
	}

	xw := &xerrorBase{xerror: new(xerror)}
	xw.Code1 = code
	xw.Msg = msg
	xw.xrr = errors.New(code)
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

	fmt.Println(handle(err, "").(*xerror).p())
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
	if isErrNil(err) {
		return nil
	}

	// 在xerror中xerrorBase就相当于errorString
	if _, ok := err.(*xerrorBase); ok {
		return err
	}

	u, ok := err.(interface{ Unwrap() error })
	if !ok {
		return err
	}
	return u.Unwrap()
}

func Is(err, target error) bool {
	if target == nil {
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
		if err = Unwrap(err); err == nil {
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
	for err != nil {
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

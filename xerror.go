package xerror

import (
	"errors"
	"fmt"
	"github.com/pubgo/xerror/internal/wrapper"
	"github.com/pubgo/xerror/xerror_util"
	"net/http"
	"os"
	"reflect"
)

type XErr interface {
	error
	Stack(indent ...bool) string
	Println() string
	String() string
}

// Combine combine multiple errors
func Combine(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var _errs xerrorCombine
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		_errs = append(_errs, handle(errs[i], ""))
	}

	if len(_errs) == 0 {
		return nil
	}

	return _errs
}

// Parse parse error to xerror
func Parse(err error) XErr {
	return handle(err, "")
}

func Fmt(format string, a ...interface{}) *xerrorBase {
	xrr := New(fmt.Sprintf(format, a...))
	xrr.Caller = xerror_util.CallerWithDepth(wrapper.CallDepth())
	return xrr
}

func New(code string, ms ...string) *xerrorBase {
	var msg string
	if len(ms) > 0 {
		msg = ms[0]
	}

	xw := &xerrorBase{}
	xw.Code = code
	xw.Msg = msg
	xw.Caller = xerror_util.CallerWithDepth(wrapper.CallDepth())

	return xw
}

func Try(fn func()) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			err2 := &xerror{}
			err2.Caller = xerror_util.CallerWithFunc(fn)

			switch err1 := _err.(type) {
			case error:
				err2.Cause1 = New(unwrap(err1).Error(), fmt.Sprintf("%+v", err1))
			default:
				err2.Cause1 = New(ErrUnknownType.Error(), fmt.Sprintf("%+v", err1))
			}
			err = err2
		}
	}()
	fn()
	return
}

func RespErr(err *error) {
	handleErr(err, recover())
}

var goroutineErrHandle = func(err XErr) {
	if isErrNil(err) {
		return
	}

	fmt.Println(err.Println())
}

func init() {
	go func() {
		for {
			select {
			case err := <-goroutineErrs:
				if goroutineErrHandle != nil {
					goroutineErrHandle(err)
				}
			}
		}
	}()
}

func SetGoroutineErrHandle(fn func(err XErr)) {
	if fn == nil {
		log.Fatal("the parameters should not be nil")
	}
	goroutineErrHandle = fn
}

var goroutineErrs = make(chan *xerror, 1)

func RespGoroutine() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}
	goroutineErrs <- handle(err, "")
}

func RespDebug() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	fmt.Println(handle(err, "").p())
	wrapper.PrintStack()
}

// Resp
func Resp(f func(err XErr)) {
	var err error
	handleErr(&err, recover())
	if err == nil {
		return
	}

	if err, ok := err.(XErr); ok {
		f(err.(XErr))
		return
	}
	f(&xerror{Cause1: err, Caller: xerror_util.CallerWithDepth(wrapper.CallDepth() + 1)})
}

func RespExit() {
	var err error
	handleErr(&err, recover())
	if isErrNil(err) {
		return
	}

	fmt.Println(handle(err, "").p())
	wrapper.PrintStack()
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
func ExitErr(dat interface{}, err error) interface{} {
	if isErrNil(err) {
		return dat
	}
	fmt.Println(handle(err, "").p())
	wrapper.PrintStack()
	os.Exit(1)
	return nil
}

// ExitF
func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}
	fmt.Println(handle(err, msg, args...).p())
	wrapper.PrintStack()
	os.Exit(1)
}

func Exit(err error) {
	if isErrNil(err) {
		return
	}
	fmt.Println(handle(err, "").p())
	wrapper.PrintStack()
	os.Exit(1)
}

// FamilyAs Check if *err belongs to *target's family
func FamilyAs(err error, target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	for err != nil {
		if x, ok := err.(interface{ FamilyAs(interface{}) bool }); ok && x.FamilyAs(target) {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

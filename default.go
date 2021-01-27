package xerror

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/pubgo/xerror/xerror_abc"
)

// PanicCombine combine multiple errors
func PanicCombine(errs ...error) {
	if len(errs) == 0 {
		return
	}

	var errs1 combine
	for i := range errs {
		if isErrNil(errs[i]) {
			continue
		}

		errs1 = append(errs1, handle(errs[i]))
	}

	if len(errs1) == 0 {
		return
	}

	panic(errs1)
}

// Parse parse error to xerror
func Parse(err error) xerror_abc.XErr {
	if isErrNil(err) {
		return nil
	}

	return handle(err)
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

func As(err error, target interface{}) bool { return FamilyAs(err, target) }

// FamilyAs Assert if *err belongs to *target's family
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

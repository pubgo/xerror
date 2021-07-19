package xerror

import (
	"fmt"
	"os"
)

// PanicErrs combine multiple errors
func PanicErrs(errs ...error) {
	if len(errs) == 0 {
		return
	}

	panic(Combine(errs...))
}

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

func PanicErr(ret interface{}, err error) interface{} {
	if isErrNil(err) {
		return ret
	}

	panic(handle(err))
}

func PanicBytes(ret []byte, err error) []byte {
	if isErrNil(err) {
		return ret
	}

	panic(handle(err))
}

func PanicStr(ret string, err error) string {
	if isErrNil(err) {
		return ret
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

func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprintf(msg, args...) }).stackString())
	printStack()
	os.Exit(1)
}

func ExitErr(dat interface{}, err error) interface{} {
	if isErrNil(err) {
		return dat
	}

	p(handle(err).stackString())
	printStack()
	os.Exit(1)
	return nil
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

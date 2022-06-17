package funk

import (
	"fmt"
	"os"
)

func Panic(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(handle(err, func(err *xerror) { err.Detail = fmt.Sprint(args...) }))
}

func PanicF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(handle(err, func(err *xerror) { err.Detail = fmt.Sprintf(msg, args...) }))
}

func PanicErr[T any](ret T, err error) T {
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

	p(handle(err, func(err *xerror) { err.Detail = fmt.Sprint(args...) }).debugString())
	printStack()
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Detail = fmt.Sprintf(msg, args...) }).debugString())
	printStack()
	os.Exit(1)
}

func ExitErr[T any](ret T, err error) T {
	if isErrNil(err) {
		return ret
	}

	p(handle(err).debugString())
	printStack()
	os.Exit(1)
	return ret
}

func Wrap(err error, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	return handle(err, func(err *xerror) { err.Detail = fmt.Sprint(args...) })
}

func WrapF(err error, msg string, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	return handle(err, func(err *xerror) { err.Detail = fmt.Sprintf(msg, args...) })
}

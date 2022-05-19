package xerror

import (
	"fmt"
	"os"
)

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

func PanicErr[A any](ret A, err error) A {
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

func ExitErr[A any](dat A, err error) A {
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

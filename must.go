package funk

import (
	"fmt"
	"os"
)

func Must(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(handle(err, func(err *xerror) { err.Detail = fmt.Sprint(args...) }))
}

func MustMsg(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(handle(err, func(err *xerror) { err.Detail = fmt.Sprintf(msg, args...) }))
}

func Must1[T any](ret T, err error) T {
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

func ExitMsg(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Detail = fmt.Sprintf(msg, args...) }).debugString())
	printStack()
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
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

func WrapFn(err error, fn func(err XErr) XErr) error {
	if isErrNil(err) {
		return nil
	}

	return fn(handle(err))
}

func WrapMsg(err error, msg string, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	return handle(err, func(err *xerror) { err.Detail = fmt.Sprintf(msg, args...) })
}

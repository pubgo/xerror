package assert

import (
	"fmt"
	"os"

	"github.com/pubgo/funk/typex"
	"github.com/pubgo/funk/xerr"
)

func Must(err error, args ...interface{}) {
	if err == nil {
		return
	}

	panic(xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprint(args...) }))
}

func MustF(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	panic(xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprintf(msg, args...) }))
}

func Must1[T any](ret T, err error) T {
	if err == nil {
		return ret
	}

	panic(xerr.WrapXErr(err))
}

func Exit(err error, args ...interface{}) {
	if err == nil {
		return
	}

	xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprint(args...) }).DebugPrint()
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprintf(msg, args...) }).DebugPrint()
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
	if err == nil {
		return ret
	}

	xerr.WrapXErr(err).DebugPrint()
	os.Exit(1)
	return ret
}

func Unwrap[T any](res T, err error) typex.Result[T] {
	return typex.OK(res, err)
}

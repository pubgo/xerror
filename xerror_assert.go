package xerror

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/kr/pretty"
)

func AssertNil(args ...interface{}) {
	for i := range args {
		if args[i] == nil || reflect.ValueOf(args[i]).IsZero() {
			panic(handle(ErrIsNil, func(err *xerror) { err.Msg = pretty.Sprint(args[i], "is nil") }))
		}
	}
}

func AssertEqual(a, b interface{}, opts ...cmp.Option) {
	if cmp.Equal(a, b, opts...) {
		return
	}

	panic(handle(ErrAssert, func(err *xerror) { err.Msg = fmt.Sprintf("[%#v] not match [%#v]", a, b) }))
}

func AssertNotEqual(a, b interface{}, opts ...cmp.Option) {
	if !cmp.Equal(a, b, opts...) {
		return
	}

	panic(handle(ErrAssert, func(err *xerror) { err.Msg = fmt.Sprintf("[%#v] match [%#v]", a, b) }))
}

func AssertErr(b bool, format string, a ...interface{}) error {
	if !b {
		return nil
	}

	return fmt.Errorf(format, a...)
}

func Assert(b bool, format string, a ...interface{}) {
	if !b {
		return
	}

	panic(handle(ErrAssert, func(err *xerror) { err.Msg = fmt.Sprintf(format, a...) }))
}

func AssertFn(b bool, fn func() string) {
	if !b {
		return
	}

	panic(handle(ErrAssert, func(err *xerror) { err.Msg = fn() }))
}

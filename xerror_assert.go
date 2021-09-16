package xerror

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

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

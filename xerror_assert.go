package xerror

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

func AssertEqual(a, b interface{}, opts ...cmp.Option) {
	if !cmp.Equal(a, b, opts...) {
		panic(handle(errors.New("match error"), func(err *xerror) {
			err.Detail = fmt.Sprintf("a=%#v b=%#v", a, b)
		}))
	}
}

func AssertNotEqual(a, b interface{}, opts ...cmp.Option) {
	if cmp.Equal(a, b, opts...) {
		panic(handle(errors.New("not match error"), func(err *xerror) {
			err.Detail = fmt.Sprintf("a=%#v b=%#v", a, b)
		}))
	}
}

func Assert(b bool, format string, a ...interface{}) {
	if b {
		panic(handle(fmt.Errorf(format, a...)))
	}
}

func AssertErr(b bool, err error) {
	if b {
		panic(handle(err))
	}
}

func AssertFn(b bool, fn func() error) {
	if b {
		panic(handle(fn()))
	}
}

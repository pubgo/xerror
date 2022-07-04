package funk

import (
	"fmt"

	"github.com/pubgo/funk/xerr"
)

func Assert(b bool, format string, a ...interface{}) {
	if b {
		panic(xerr.WrapXErr(fmt.Errorf(format, a...)))
	}
}

func AssertErr(b bool, err error) {
	if b {
		panic(xerr.WrapXErr(err))
	}
}

func AssertFn(b bool, fn func() error) {
	if b {
		panic(xerr.WrapXErr(fn()))
	}
}

package funk

import (
	"fmt"
)

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

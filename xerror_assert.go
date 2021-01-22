package xerror

import (
	"fmt"
)

func Assert(b bool, format string, a ...interface{}) {
	if !b {
		return
	}

	Panic(fmt.Errorf(format, a...))
}

func AssertFn(b bool, fn func() error) {
	if !b {
		return
	}

	Panic(fn())
}

func If(a bool, b interface{}, c interface{}) interface{} {
	if a {
		return b
	}
	return c
}

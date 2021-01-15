package xerror

import (
	"fmt"
	"reflect"
)

func AssertNil(val interface{}, format string, a ...interface{}) {
	if val != nil && !reflect.ValueOf(val).IsNil() {
		return
	}
	Next().Panic(fmt.Errorf(format, a...))
}

func Assert(b bool, format string, a ...interface{}) {
	if !b {
		return
	}
	Next().Panic(fmt.Errorf(format, a...))
}

func If(a bool, b interface{}, c interface{}) interface{} {
	if a {
		return b
	}
	return c
}

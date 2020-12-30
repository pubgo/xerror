package xerror

import (
	"fmt"
	"reflect"

	"github.com/pubgo/xerror/xerror_util"
)

func AssertNotNil(val interface{}, a ...interface{}) {
	if val == nil || reflect.ValueOf(val).IsZero() {
		Next().PanicF(New("[val] is nil"), handleArgs(a...))
	}
}

func Assert(valid bool, a ...interface{}) {
	if valid {
		return
	}

	Next().PanicF(New("[valid] is false"), handleArgs(a...))
}

func handleArgs(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}

	switch reflect.TypeOf(a[0]).Kind() {
	case reflect.Func:
		var ret string
		xerror_util.Func(a[0])(a[1:]...)(func(s string) { ret = s })
		return ret
	case reflect.String:
		return fmt.Sprintf(a[0].(string), a[1:]...)
	default:
		panic("[a] type error")
	}
}

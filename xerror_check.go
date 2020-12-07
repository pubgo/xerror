package xerror

import (
	"errors"
	"reflect"
)

func CheckNil(val interface{}) {
	if val == nil {
		With(WithCaller(1)).Panic(errors.New("[val] is nil"))
	}

	vf := reflect.ValueOf(val)

	switch vf.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if vf.IsZero() {
			With(WithCaller(1)).Panic(errors.New("[val] is nil"))
		}
	}
}

func Check(valid bool) {
	if valid {
		return
	}

	With(WithCaller(1)).Panic(errors.New("[valid] is false"))
}

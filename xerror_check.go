package xerror

import (
	"errors"
	"reflect"
)

func CheckNil(val interface{}, format string, a ...interface{}) {
	if val == nil {
		Next().ExitF(errors.New("[val] is nil"), format, a...)
	}

	vf := reflect.ValueOf(val)
	switch vf.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if vf.IsZero() {
			Next().ExitF(errors.New("[val] is nil"), format, a...)
		}
	}
}

func Check(valid bool, format string, a ...interface{}) {
	if !valid {
		return
	}

	Next().ExitF(errors.New("[valid] is false"), format, a...)
}

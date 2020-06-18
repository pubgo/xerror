package xerror

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

func handleErr(err *error, _err interface{}) {
	if _err == nil {
		return
	}

	switch _err := _err.(type) {
	case *xerror:
		*err = _err
	case *xerrorBase:
		*err = _err.xerror
	case error:
		err1 := getXerror()
		err1.xrr = _err
		err1.Msg = fmt.Sprintf("%#v", _err)
		*err = err1
	case string:
		err1 := getXerror()
		err1.xrr = errors.New(_err)
		*err = err1
	default:
		err1 := getXerror()
		err1.xrr = ErrUnknownType
		err1.Msg = fmt.Sprintf("%#v", _err)
		*err = err1
	}
}

func handle(err error, msg string, args ...interface{}) error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	err2 := getXerror()
	err2.Msg = msg
	err2.Caller = callerWithDepth(callDepth + 1)

	switch err := err.(type) {
	case *xerrorBase:
		err2.xrr = err
		err2.Sub = err.xerror
		err2.Code1 = err.Code1
	case *xerror:
		err2.Sub = err
		err2.xrr = err.xrr
		err2.Code1 = err.Code1

		err.xrr = nil
		err.Code1 = ""
	default:
		err2.xrr = err
	}

	return err2
}

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func callerWithDepth(callDepths ...int) string {
	var cd = callDepth
	if len(callDepths) > 0 {
		cd = callDepths[0]
	}

	var pcs = make([]uintptr, 1)
	if runtime.Callers(cd, pcs[:]) == 0 {
		return ""
	}

	f := frame(pcs[0])
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown type"
	}

	file, line := fn.FileLine(f.pc())
	return file + ":" + strconv.Itoa(line)
}

func callerWithFunc(fn reflect.Value) string {
	if !fn.IsValid() || fn.IsNil() || fn.Kind() != reflect.Func {
		panic(ErrNotFuncType)
	}
	var _fn = fn.Pointer()
	var file, line = runtime.FuncForPC(_fn).FileLine(_fn)
	return file + ":" + strconv.Itoa(line)
}

func isErrNil(err error) bool {
	return err == nil || err == ErrDone
}

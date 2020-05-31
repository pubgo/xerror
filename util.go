package xerror

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func handleErr(err *error, _err interface{}) {
	if _err == nil {
		return
	}

	switch _err.(type) {
	case *xerror:
		*err = _err.(*xerror)
	case error:
		err1 := getXerror()
		err1.xrr = _err.(error)
		*err = err1
	case string:
		err1 := getXerror()
		err1.xrr = errors.New(_err.(string))
		*err = err1
	default:
		logger.Fatalf("unknown type, %#v", _err)
	}
}

func handle(err error, msg string, args ...interface{}) error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	err2 := getXerror()
	err2.Msg = msg
	err2.Caller = callerWithDepth()
	if err1, ok := err.(*xerror); ok {
		err2.Sub = err1
		err2.xrr = err1.xrr
		err1.xrr = nil

		err2.code = err1.code
		err1.code = 0
	} else {
		err2.xrr = err
	}

	return err2
}

type Frame uintptr

func (f Frame) pc() uintptr { return uintptr(f) - 1 }

func callerWithDepth(callDepths ...int) string {
	var cd = callDepth
	if len(callDepths) > 0 {
		cd = callDepths[0]
	}

	var pcs = make([]uintptr, 1)
	if runtime.Callers(cd, pcs[:]) == 0 {
		return ""
	}

	f := Frame(pcs[0])
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown func"
	}

	file, line := fn.FileLine(f.pc())
	return file + ":" + strconv.Itoa(line)
}

func callerWithFunc(fn reflect.Value) string {
	if !fn.IsValid() || fn.IsNil() || fn.Kind() != reflect.Func {
		logger.Fatal("func error")
	}
	var _fn = fn.Pointer()
	var file, line = runtime.FuncForPC(_fn).FileLine(_fn)
	return file + ":" + strconv.Itoa(line)
}

func isErrNil(err error) bool {
	return err == nil || err == ErrDone
}

func env(es ...string) string {
	for _, e := range es {
		if v := os.Getenv(strings.ToUpper(e)); v != "" {
			return v
		}
	}
	return ""
}

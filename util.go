package xerror

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"

	"github.com/pubgo/xerror/internal/wrapper"
)

func handleErr(err *error, _err interface{}) {
	if _err == nil {
		return
	}

	switch _err := _err.(type) {
	case error:
		*err = _err
	case string:
		*err = errors.New(_err)
	default:
		*err = New("unknown type", fmt.Sprintf("%+v", _err))
	}
}

func handle(err error, msg string, args ...interface{}) *xerror {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	err2 := &xerror{}
	err2.Msg = msg
	err2.Caller = callerWithDepth(wrapper.CallDepth() + 1)
	err2.Cause1 = err
	return err2
}

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func callerWithDepth(callDepths ...int) string {
	if !wrapper.IsCaller() {
		return ""
	}

	var cd = wrapper.CallDepth()
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
		log.Fatalln("not func type")
	}
	var _fn = fn.Pointer()
	var file, line = runtime.FuncForPC(_fn).FileLine(_fn)
	return file + ":" + strconv.Itoa(line)
}

func isErrNil(err error) bool {
	return err == nil || err == ErrDone
}

func trans(err error) *xerror {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *xerrorBase:
		return &xerror{
			Msg:    err.Msg,
			Caller: err.Caller,
		}
	case *xerror:
		return err
	default:
		return nil
	}
}

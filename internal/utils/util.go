package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/pubgo/xerror/xerror_conf"
)

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func CallerWithDepth(cd int) string {
	if !xerror_conf.Conf.EnableCaller {
		return ""
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
	return fmt.Sprintf("%s:%d", file, line)
}

func CallerWithFunc(fn interface{}) string {
	if fn == nil {
		panic("[fn] is nil")
	}

	if !xerror_conf.Conf.EnableCaller {
		return ""
	}

	var _fn = reflect.ValueOf(fn)
	if !_fn.IsValid() || _fn.Kind() != reflect.Func || _fn.IsNil() {
		panic("[fn] is not func type or type is nil")
	}

	var _e = runtime.FuncForPC(_fn.Pointer())
	var file, line = _e.FileLine(_fn.Pointer())

	ma := strings.Split(_e.Name(), ".")
	return fmt.Sprintf("%s:%d %s", file, line, ma[len(ma)-1])
}

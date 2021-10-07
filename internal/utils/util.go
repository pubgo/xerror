package utils

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/valyala/bytebufferpool"
)

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func CallerWithDepth(cd int) string {
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

func CallerWithFunc(fn interface{}) string {
	if fn == nil {
		panic("[fn] is nil")
	}

	var _fn = reflect.ValueOf(fn)
	if !_fn.IsValid() || _fn.Kind() != reflect.Func || _fn.IsNil() {
		panic("not func type or type is nil")
	}

	var _e = runtime.FuncForPC(_fn.Pointer())
	var file, line = _e.FileLine(_fn.Pointer())

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")

	ma := strings.Split(_e.Name(), ".")
	buf.WriteString(ma[len(ma)-1])
	return buf.String()
}

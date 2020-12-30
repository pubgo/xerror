package xerror_util

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/pubgo/xerror/xerror_envs"
)

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func CallerWithDepth(callDepths ...int) string {
	if !xerror_envs.IsCallerVal() {
		return ""
	}

	var cd = xerror_envs.CallDepthVal()
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

func CallerWithFunc(fn interface{}) string {
	if fn == nil {
		log.Fatalln("params is nil")
	}

	var _fn = reflect.ValueOf(fn)
	if !_fn.IsValid() || _fn.IsNil() || _fn.Kind() != reflect.Func {
		log.Fatalln("not func type or type is nil")
	}

	var _e = runtime.FuncForPC(_fn.Pointer())
	var file, line = _e.FileLine(_fn.Pointer())

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")

	ma := strings.Split(_e.Name(), ".")
	buf.WriteString(ma[len(ma)-1])
	return buf.String()
}

func FuncRaw(fn interface{}) func(...interface{}) []reflect.Value {
	if fn == nil {
		panic("[fn] is nil")
	}

	vfn, ok := fn.(reflect.Value)
	if !ok {
		vfn = reflect.ValueOf(fn)
	}
	if !vfn.IsValid() || vfn.Kind() != reflect.Func || vfn.IsNil() {
		panic("[fn] type error or nil")
	}

	var tfn = vfn.Type()
	var numIn = tfn.NumIn()
	var variadicType reflect.Type
	if tfn.IsVariadic() {
		variadicType = tfn.In(numIn - 1)
	}

	return func(args ...interface{}) []reflect.Value {
		if variadicType == nil && numIn != len(args) || variadicType != nil && len(args) < numIn-1 {
			panic(fmt.Sprintf("the input params of func is not match, func: %s, numIn:%d numArgs:%d", tfn, numIn, len(args)))
		}

		var _args = valueGet()
		for _, k := range args {
			var vk reflect.Value
			if k == nil {
				vk = reflect.ValueOf(k)
			} else if k1, ok := k.(reflect.Value); ok {
				vk = k1
			} else {
				vk = reflect.ValueOf(k)
			}
			_args = append(_args, vk)
		}

		for i, k := range _args {
			if i >= numIn {
				if variadicType == nil {
					panic(fmt.Sprintf("[variadicType] should not be nil, args:%s, fn:%s", fmt.Sprint(args...), tfn))
				}

				_args[i] = reflect.Zero(variadicType)
				continue
			}

			if !k.IsValid() {
				_args[i] = reflect.Zero(tfn.In(i))
			}
		}

		defer func() {
			valuePut(_args)
			if err := recover(); err != nil {
				panic(fmt.Sprintf("[vfn.Call] panic, err:%#v, args:%s, fn:%s", err, valueStr(_args...), tfn))
			}
		}()
		return vfn.Call(_args)
	}
}

func Func(fn interface{}) func(...interface{}) func(...interface{}) {
	vfn := FuncRaw(fn)
	return func(args ...interface{}) func(...interface{}) {
		ret := vfn(args...)
		return func(fns ...interface{}) {
			if len(fns) == 0 {
				return
			}

			if fns[0] == nil {
				panic("[fns] is nil")
			}

			cfn, ok := fns[0].(reflect.Value)
			if !ok {
				cfn = reflect.ValueOf(fns[0])
			}
			if !cfn.IsValid() || cfn.Kind() != reflect.Func || cfn.IsNil() {
				panic("[fns] type error or nil")
			}

			tfn := reflect.TypeOf(fn)
			if cfn.Type().NumIn() != tfn.NumOut() {
				panic(fmt.Sprintf("the input num and output num of the callback func is not match, [%d]<->[%d]",
					cfn.Type().NumIn(), tfn.NumOut()))
			}

			if cfn.Type().NumIn() != 0 && cfn.Type().In(0) != tfn.Out(0) {
				panic(fmt.Sprintf("the output type of the callback func is not match, [%s]<->[%s]",
					cfn.Type().In(0), tfn.Out(0)))
			}

			defer func() {
				valuePut(ret)
				if err := recover(); err != nil {
					panic(fmt.Sprintf("[cfn.Call] panic, err:%#v, args:%s, fn:%s", err, valueStr(ret...), cfn.Type()))
				}
			}()

			cfn.Call(ret)
		}
	}
}

var _valuePool = sync.Pool{
	New: func() interface{} {
		return []reflect.Value{}
	},
}

func valueGet() []reflect.Value {
	return _valuePool.Get().([]reflect.Value)
}

func valuePut(v []reflect.Value) {
	_valuePool.Put(v[:0])
}

func valueStr(values ...reflect.Value) string {
	var data []interface{}
	for _, dt := range values {
		var val interface{} = nil
		if dt.IsValid() {
			val = dt.Interface()
		}
		data = append(data, val)
	}
	return fmt.Sprint(data...)
}

package xerror_util

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
		log.Fatalln("params is nil")
	}

	var _fn = reflect.ValueOf(fn)
	if !_fn.IsValid() || _fn.Kind() != reflect.Func || _fn.IsNil() {
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

func FuncValue(fn interface{}) func(...reflect.Value) []reflect.Value {
	if fn == nil {
		log.Fatalln("[fn] is nil")
	}

	vfn, ok := fn.(reflect.Value)
	if !ok {
		vfn = reflect.ValueOf(fn)
	}
	if !vfn.IsValid() || vfn.Kind() != reflect.Func || vfn.IsNil() {
		log.Fatalln("[fn] type error or nil")
	}

	var tfn = vfn.Type()
	var numIn = tfn.NumIn()
	var variadicType reflect.Type
	if tfn.IsVariadic() {
		variadicType = tfn.In(numIn - 1)
	}

	return func(args ...reflect.Value) []reflect.Value {
		if variadicType == nil && numIn != len(args) || variadicType != nil && len(args) < numIn-1 {
			log.Fatalf("the input params of func is not match, func: %s, numIn:%d numArgs:%d\n", tfn, numIn, len(args))
		}

		for i := range args {
			if i >= numIn {
				if variadicType == nil {
					log.Fatalf("[variadicType] should not be nil, args:%s, fn:%s", valueStr(args...), tfn)
				}

				args[i] = reflect.Zero(variadicType)
				continue
			}

			if !args[i].IsValid() {
				args[i] = reflect.Zero(tfn.In(i))
			}
		}

		defer func() {
			valuePut(args)
			if err := recover(); err != nil {
				panic(fmt.Sprintf("[vfn.Call] panic, err:%#v, args:%s, fn:%s", err, valueStr(args...), tfn))
			}
		}()
		return vfn.Call(args)
	}
}

func FuncRaw(fn interface{}) func(...interface{}) []reflect.Value {
	vfn := FuncValue(fn)
	return func(args ...interface{}) []reflect.Value {
		var _args = valueGet()
		for i := range args {
			var vk reflect.Value
			if args[i] == nil {
				vk = reflect.ValueOf(args[i])
			} else if k1, ok := args[i].(reflect.Value); ok {
				vk = k1
			} else {
				vk = reflect.ValueOf(args[i])
			}
			_args = append(_args, vk)
		}
		return vfn(_args...)
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
				log.Fatalln("[fns] is nil")
			}

			cfn, ok := fns[0].(reflect.Value)
			if !ok {
				cfn = reflect.ValueOf(fns[0])
			}
			if !cfn.IsValid() || cfn.Kind() != reflect.Func || cfn.IsNil() {
				log.Fatalln("[fns] type error or nil")
			}

			tfn := reflect.TypeOf(fn)
			if cfn.Type().NumIn() != tfn.NumOut() {
				log.Fatalf("the input num and output num of the callback func is not match, [%d]<->[%d]\n",
					cfn.Type().NumIn(), tfn.NumOut())
			}

			if cfn.Type().NumIn() != 0 && cfn.Type().In(0) != tfn.Out(0) {
				log.Fatalf("the output type of the callback func is not match, [%s]<->[%s]\n",
					cfn.Type().In(0), tfn.Out(0))
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

var valuePool = sync.Pool{
	New: func() interface{} {
		return make([]reflect.Value, 0, 1)
	},
}

func valueGet() []reflect.Value {
	return valuePool.Get().([]reflect.Value)
}

func valuePut(v []reflect.Value) {
	valuePool.Put(v[:0])
}

func valueStr(values ...reflect.Value) string {
	var data []interface{}
	for i := range values {
		var val interface{} = nil
		if values[i].IsValid() {
			val = values[i].Interface()
		}
		data = append(data, val)
	}
	return fmt.Sprint(data...)
}

package xerror_test

import (
	"fmt"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/errs"
	"github.com/pubgo/xtest"
	"testing"
)

func init22(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)

	//fmt.Println(a...)
	//xrr.Panic(fmt.Errorf("ss"))
	//Exit(New(""))
	//_ = fmt.Sprintf("ss")
	//_ = fmt.Errorf("ss")
	//_ = "ss" + "sss"
	//xrr.Panic(nil)
	//xerror.PanicF(nil, "sssss %#v", a)
	xerror.PanicF(errs.ErrBadRequest, "ssssss wrap")
	//xerror.PanicF(fmt.Errorf("ss"), "sssss %#v", a)
	return
}

func init21(a ...interface{}) (err error) {
	//defer xerror.RespErr(&err)
	defer xerror.Resp(func(_err xerror.XRErr) {
		_ = _err.Error()
		//fmt.Println(_err.Error(), _err.Code())
	})

	//fmt.Println(a...)
	//xrr.Panic(fmt.Errorf("ss"))
	//xrr.PanicF(init22(a...), "sssss %#v", a)
	xerror.Panic(init22(a...))
	return
}

func TestName(t *testing.T) {

	fmt.Println(init21(1, 2, 3))
	//Exit(init21(1, 2, 3))
}

func TestTry(t *testing.T) {
	fmt.Println(xerror.Try(func() {
		panic("hello")
	}))
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		init21(1, 2, 3)
	}
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			defer xerror.Resp(func(_err xerror.XRErr) {
				_err.Error()
			})

			xerror.PanicF(nil, "hello")
		}()
	}
}

func TestNew(t *testing.T) {
	fn := xtest.TestFuncWith(func(code string, ms ...string) {
		defer xerror.RespExit()
		xrr := xerror.New(code, ms...)
		if code != xrr.Code() {
			xerror.Exit(xrr)
		}
	})
	fn.In("错误信息的简介和标志, 类似于404", "", xtest.RangeString(10, 100))
	fn.In("错误信息的介绍", "")
	fn.In("错误信息的介绍", "11")
	fn.Do()
}

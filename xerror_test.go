package xerror_test

import (
	"encoding/json"
	"fmt"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/errs"
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
	defer xerror.RespErr(&err)
	//defer xerror.Resp(func(_err xerror.XErr) {
		//fmt.Println(_err.Stack())
		//_ = _err.Error()
		//fmt.Println(_err.Error(), _err.Code())
	//})

	//fmt.Println(a...)
	//xrr.Panic(fmt.Errorf("ss"))
	//xrr.PanicF(init22(a...), "sssss %#v", a)
	xerror.Panic(init22(a...))
	return
}

func TestName(t *testing.T) {
	sss:=init21(1, 2, 3)
	dt, _ := json.Marshal(sss)
	fmt.Println( string(dt))
	xerror.Exit(sss)

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
			defer xerror.Resp(func(_err xerror.XErr) {
				_err.Error()
			})

			xerror.PanicF(nil, "hello")
		}()
	}
}

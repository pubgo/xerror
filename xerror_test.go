package xerror_test

import (
	"fmt"
	"github.com/pubgo/xerror"
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
	xerror.Panic(xerror.ErrBadRequest.Wrap(fmt.Errorf("ssssss wrap")))
	//xerror.PanicF(fmt.Errorf("ss"), "sssss %#v", a)
	return
}

func init21(a ...interface{}) (err error) {
	//defer xerror.RespErr(&err)
	defer xerror.Resp(func(_err xerror.XErr) {
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
	xerror.Debug = true

	fmt.Println(init21(1, 2, 3))
	//Exit(init21(1, 2, 3))
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
				//fmt.Println(_err.Error(), _err.Code())
				_err.Error()
			})

			xerror.PanicF(nil, "hello")
		}()
	}
}

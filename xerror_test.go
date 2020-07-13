package xerror_test

import (
	"fmt"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/errs"
	"github.com/pubgo/xerror/xerror_core"
	"testing"
)

func a1(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.PanicF(errs.ErrBadRequest, "test %+v", a)
	return
}

func a2(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.Panic(a1(a...))
	return
}

func TestName(t *testing.T) {
	xerror_core.IsCaller = false
	defer xerror.Resp(func(err xerror.XErr) {
		fmt.Println(err.Stack())
	})
	xerror.Panic(a2(1, 2, 4, 5))
}

func TestTry(t *testing.T) {
	fmt.Println(xerror.Try(func() {
		panic("hello")
	}))
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.PanicF(errs.ErrBadRequest, "测试Panic")
			return
		}()
	}
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.PanicF(nil, "测试NoPanic")
			return
		}()
	}
}

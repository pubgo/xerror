package xerror_test

import (
	"fmt"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/xerror_http"
	"log"
	"testing"
)

func check(b bool) {
	if !b {
		log.Fatalln("")
	}
}

func a1(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.PanicF(xerror_http.ErrBadRequest, "test %+v", a)
	return
}

func a2(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.Panic(a1(a...))
	return
}

func TestName(t *testing.T) {
	defer xerror.Resp(func(err xerror.XErr) {
		fmt.Println(err.Stack(true))
	})
	xerror.Panic(a2(1, 2, 4, 5))
}

func TestAs(t *testing.T) {
	check(xerror.FamilyAs(a2(1, 2, 4, 5), xerror_http.ErrHttp)==true)
	check(xerror.FamilyAs(a2(1, 2, 4, 5), xerror_http.ErrBadRequest)==true)
	check(xerror.FamilyAs(a2(1, 2, 4, 5), xerror_http.ErrNotFound)==false)
}

func TestExit(t *testing.T) {
	xerror.Exit(a2(1, 2, 4, 5))
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
			xerror.PanicF(xerror_http.ErrBadRequest, "测试Panic")
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

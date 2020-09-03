package xerror_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/xerror_core"
	"github.com/pubgo/xerror/xerror_http"
	"github.com/pubgo/xlog"
)

func check(b bool) {
	if !b {
		log.Fatalln("")
	}
}

func panic1(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.PanicF(xerror_http.ErrBadRequest, "panic1 %+v", a)
	return
}

func panic2(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.PanicF(panic1(a...), "panic2 %+v", a)
	return
}

func panicWrap(a ...interface{}) (err error) {
	return xerror.WrapF(panic2(a...), "panicWrap %+v", a)
}

func TestCombine(t *testing.T) {
	defer xerror.Resp(func(err xerror.XErr) {
		fmt.Println(err.Stack(true))
	})
	xerror.Panic(xerror.Combine(panicWrap(1, 2, 4, 5), panicWrap(1, 2, 4, 5)))
}

func TestLog(t *testing.T) {
	xlog.InfoF("gg \n%v", xerror.Parse(xerror.New("ddd", "dddnjnjnj")).Stack(true))
}

func TestStack(t *testing.T) {
	defer xerror.Resp(func(err xerror.XErr) {
		fmt.Println(err.Stack(true))
	})
	xerror.Exit(panicWrap(1, 2, 4, 5))
}

func TestAs(t *testing.T) {
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), xerror_http.ErrHttp) == true)
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), xerror_http.ErrBadRequest) == true)
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), xerror_http.ErrNotFound) == false)
}

func TestExit(t *testing.T) {
	xerror_core.PrintStack = false
	xerror.Exit(panicWrap(1, 2, 4, 5))
}

func TestTry(t *testing.T) {
	fmt.Println(xerror.Try(func() {
		panic("hello")
	}))
}

func TestRespGoroutine(t *testing.T) {
	xerror.Exit(xerror.SetGoroutineErrHandler("test", func(err xerror.XErr) {
		fmt.Println(err.Stack(true))
	}))

	go func() {
		defer xerror.RespGoroutine("test")
		xerror.Panic(panicWrap(1, 2, 4, 5))
	}()

	go func() {
		defer xerror.RespGoroutine()
		xerror.Panic(panicWrap(1, 2, 4, 5))
	}()
	time.Sleep(time.Second)
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(xerror_http.ErrBadRequest)
			return
		}()
	}
}

func BenchmarkPanicWithoutCaller(b *testing.B) {
	xerror_core.IsCaller = false
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(xerror_http.ErrBadRequest)
			return
		}()
	}
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(nil)
			return
		}()
	}
}

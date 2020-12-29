package xerror_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/internal/envs"
	"github.com/pubgo/xerror/xerror_http"
)

func TestRespNext(t *testing.T) {
	defer xerror.RespDebug()
	TestPanic1(t)
}

func TestPanic(t *testing.T) {
	defer xerror.RespJson()
	xerror.Panic(xerror.New("ok"))
}

func TestPanic1(t *testing.T) {
	defer xerror.RespRaise(func(err xerror.XErr) error {
		return xerror.WrapF(err, "test raise")
	})

	xerror.Panic(xerror.New("ok"))
}

func TestPanicWith(t *testing.T) {
	defer xerror.RespJson()

	xerror.With().Panic(xerror.New("ok"))
}

func init1Next() (err error) {
	defer xerror.RespErr(&err)
	xerror.Next().Panic(fmt.Errorf("test next"))
	return nil
}

func TestDone(t *testing.T) {
	defer xerror.RespJson()
	xerror.Done()
}

func TestNext(t *testing.T) {
	defer xerror.RespJson()
	xerror.Next().Panic(init1Next())
}

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
	//xerror.With().Panic(xerror.Combine(panicWrap(1, 2, 4, 5), panicWrap(1, 2, 4, 5)))
	xerror.Panic(xerror.Combine(panicWrap(1, 2, 4, 5), panicWrap(1, 2, 4, 5)))
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
	envs.PrintStack = false
	xerror.Exit(panicWrap(1, 2, 4, 5))
	//fmt.Printf("%s\n",panicWrap(1, 2, 4, 5))
	//fmt.Printf("%v\n",panicWrap(1, 2, 4, 5))
	//fmt.Printf("%+v\n",panicWrap(1, 2, 4, 5))
	//fmt.Printf("%#v\n\n\n",panicWrap(1, 2, 4, 5))
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
	envs.IsCaller = false
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

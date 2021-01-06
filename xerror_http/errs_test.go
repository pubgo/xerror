package xerror_http

import (
	"fmt"
	"log"
	"testing"

	"github.com/pubgo/xerror"
	"github.com/pubgo/xerror/xerror_core"
)

func panic1(a ...interface{}) (err error) {
	defer xerror.RespErr(&err)
	xerror.PanicF(ErrBadRequest, "panic1 %+v", a)
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
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), ErrHttp) == true)
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), ErrBadRequest) == true)
	check(xerror.FamilyAs(panicWrap(1, 2, 4, 5), ErrNotFound) == false)
}

func TestExit(t *testing.T) {
	xerror_core.Conf.PrintStack = false
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

func check(b bool) {
	if !b {
		log.Fatalln("")
	}
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(ErrBadRequest)
			return
		}()
	}
}

func BenchmarkPanicWithoutCaller(b *testing.B) {
	xerror_core.Conf.IsCaller = false
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RespErr(&err)
			xerror.Panic(ErrBadRequest)
			return
		}()
	}
}

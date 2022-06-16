package xerror_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/xerror"
)

func panicErr() (*xerror.Err, error) {
	return nil, fmt.Errorf("error")
}

func panicNoErr() (*xerror.Err, error) {
	return &xerror.Err{Msg: "ok"}, nil
}

func TestPanicErr(t *testing.T) {
	defer xerror.RecoverTest(t)
	var err = xerror.Try(func() {
		var ret = xerror.PanicErr(panicErr())
		fmt.Println(ret == nil)
	})
	xerror.Assert(err == nil, "failed")

	err = xerror.Try(func() {
		var ret = xerror.PanicErr(panicNoErr())
		fmt.Println(ret.Msg)
	})
	xerror.Assert(err != nil, "failed")
}

func TestRespTest(t *testing.T) {
	defer xerror.RecoverTest(t)
	testPanic1(t)
}

func TestRespNext(t *testing.T) {
	defer xerror.RecoverAndExit()
	testPanic1(t)
}

func testPanic1(t *testing.T) {
	defer xerror.RecoverAndRaise()

	//xerror.Panic(xerror.New("ok"))
	xerror.Panic(init1Next())
}

func init1Next() (err error) {
	defer xerror.RecoverErr(&err)
	xerror.Panic(fmt.Errorf("test next"))
	return nil
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer xerror.RecoverErr(&err)
			xerror.Panic(nil)
			return
		}()
	}
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				recover()
			}()

			panic("hello")
		}()
	}
}

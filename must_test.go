package funk_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk"
	"github.com/stretchr/testify/assert"
)

func panicErr() (*funk.Err, error) {
	return nil, fmt.Errorf("error")
}

func panicNoErr() (*funk.Err, error) {
	return &funk.Err{Msg: "ok"}, nil
}

func TestPanicErr(t *testing.T) {
	var is = assert.New(t)
	is.Panics(func() {
		var ret = funk.Must1(panicErr())
		fmt.Println(ret == nil)
	})

	is.NotPanics(func() {
		var ret = funk.Must1(panicNoErr())
		fmt.Println(ret.Msg)
	})
}

func TestRespTest(t *testing.T) {
	testPanic1(t)
}

func TestRespNext(t *testing.T) {
	defer funk.RecoverAndExit()
	testPanic1(t)
}

func testPanic1(t *testing.T) {
	defer funk.RecoverAndRaise()

	//Xerror.Must(Xerror.New("ok"))
	funk.Must(init1Next())
}

func init1Next() (err error) {
	defer funk.RecoverErr(&err)
	funk.Must(fmt.Errorf("test next"))
	return nil
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			defer funk.RecoverErr(&err)
			funk.Must(nil)
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

package xerror_test

import (
	"fmt"
	"testing"

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
	defer xerror.RespExit()
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

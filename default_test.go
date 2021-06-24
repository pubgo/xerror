package xerror_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/xerror"
)

func TestParseWith(t *testing.T) {
	var err = fmt.Errorf("hello error")
	xerror.ParseWith(err, func(err error) {
		fmt.Printf("%v\n", err)
	})
}

func TestRespTest(t *testing.T) {
	defer xerror.RespTest(t)
	TestPanic1(t)
}

func TestRespNext(t *testing.T) {
	defer xerror.RespExit("TestRespNext")
	TestPanic1(t)
}

func TestPanic1(t *testing.T) {
	//defer xerror.RespExit()
	defer xerror.RespRaise(func(err xerror.XErr) error {
		return xerror.WrapF(err, "test raise")
	})

	//xerror.Panic(xerror.New("ok"))
	xerror.Panic(fmt.Errorf("ss"))
}

func init1Next() (err error) {
	defer xerror.RespErr(&err)
	xerror.Panic(fmt.Errorf("test next"))
	return nil
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

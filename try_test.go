package funk

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/xerr"
)

func testFunc() (err error) {
	defer RecoverErr(&err, func(err xerr.XErr) xerr.XErr {
		return err.WrapF("test func")
	})
	Must(Err{Msg: "test error"})
	return
}

func TestTryLog(t *testing.T) {
	TryAndLog(func() {
		Must(testFunc())
	})
}

func TestTryCatch(t *testing.T) {
	TryCatch(
		func() { panic("ok") },
		func(err xerr.XErr) {
			fmt.Println(err.Error(), err)
		})
}

func TestTryThrow(t *testing.T) {
	TryThrow(func() {
		panic("abc")
	})
}

func TestTryVal(t *testing.T) {
	var v = TryRet(func() (*Err, error) {
		return &Err{Msg: "ok"}, nil
	}, func(err xerr.XErr) {
		fmt.Println(err)
	})
	fmt.Println(v)
}

package funk

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/funkonf"
)

func testFunc() (err error) {
	defer RecoverErr(&err, func(err XErr) XErr {
		return err.WrapF("test func")
	})
	Must(Err{Msg: "test error"})
	return
}

func TestTryLog(t *testing.T) {
	funkonf.Conf.Debug = true
	TryAndLog(func() {
		Must(testFunc())
	})
}

func TestTryCatch(t *testing.T) {
	TryCatch(
		func() { panic("ok") },
		func(err error) {
			fmt.Println(err.Error(), err)
		})
}

func TestTryThrow(t *testing.T) {
	defer RecoverTest(t)

	TryThrow(func() {
		panic("abc")
	})
}

func TestTryVal(t *testing.T) {
	defer RecoverTest(t)

	var v = TryRet(func() (*Err, error) {
		return &Err{Msg: "ok"}, nil
	}, func(err error) {
		fmt.Println(err)
	})
	fmt.Println(v)
}

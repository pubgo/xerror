package funk

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/xerr"
)

func testFunc() (err error) {
	defer recovery.Err(&err, func(err xerr.XErr) xerr.XErr {
		return err.WrapF("test func")
	})
	assert.Must(xerr.Err{Msg: "test error"})
	return
}

func TestTryLog(t *testing.T) {
	TryAndLog(func() {
		assert.Must(testFunc())
	})
}

func TestTryCatch(t *testing.T) {
	TryCatch(
		func() error { panic("ok"); return nil },
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
	var v = Try1(func() (*xerr.Err, error) {
		return &xerr.Err{Msg: "ok"}, nil
	}, func(err xerr.XErr) {
		fmt.Println(err)
	})
	fmt.Println(v)
}

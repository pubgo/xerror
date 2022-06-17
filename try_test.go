package funk

import (
	"fmt"
	"testing"
)

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

	var v = TryVal(func() (*Err, error) {
		return &Err{Msg: "ok"}, nil
	}, func(err error) {
		fmt.Println(err)
	})
	fmt.Println(v)
}

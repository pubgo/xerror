package xerror

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

}

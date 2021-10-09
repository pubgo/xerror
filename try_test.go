package xerror

import (
	"fmt"
	"testing"
)

func TestTryThrow(t *testing.T) {
	defer RespTest(t)

	TryThrow(func() {
		panic("abc")
	}, "test try throw")
}

func TestTryVal(t *testing.T) {
	defer RespTest(t)

	fmt.Println(TryVal(func() interface{} {
		return "hello"
	}))

	var a = func() {
		panic("hello")
	}

	fmt.Println(TryVal(func() interface{} {
		fmt.Println("hello")

		a()
		return nil
	}))

}

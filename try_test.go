package xerror

import "testing"

func TestTryThrow(t *testing.T) {
	defer RespTest(t)

	TryThrow(func() {
		panic("abc")
	},"test try throw")
}
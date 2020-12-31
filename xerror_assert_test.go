package xerror

import (
	"testing"
)

func TestCheckNil(t *testing.T) {
	defer RespExit()

	var a *int
	AssertNil(a, func() string {
		return "ok"
	})
}

func TestCheck(t *testing.T) {
	defer RespDebug()

	Assert(true, "aaaa")
	Assert(false, "aaaa")
}

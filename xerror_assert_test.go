package xerror

import (
	"testing"
)

func TestCheckNil(t *testing.T) {
	defer RespExit()

	var a *int
	AssertNotNil(a, func() string {
		return "ok"
	})
	//AssertNotNil(a, "aaaa")
}

func TestCheck(t *testing.T) {
	defer RespDebug()

	Assert(true, "aaaa")
	Assert(false, "aaaa")
}

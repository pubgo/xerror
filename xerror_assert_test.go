package xerror

import (
	"testing"
)

func TestCheckNil(t *testing.T) {
	defer RecoverTest(t)

	var a *int
	Assert(a == nil, "ok")
}

func try(fn func()) (err error) {
	defer RecoverErr(&err)

	fn()
	return nil
}

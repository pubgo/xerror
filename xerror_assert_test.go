package xerror

import (
	"testing"
)

func TestAssertEqual(t *testing.T) {
	defer RecoverTest(t)

	AssertEqual("hello", 1)
}

func TestCheckNil(t *testing.T) {
	defer RecoverTest(t)

	var a *int
	Assert(a == nil, "ok")
}

func TestCheck(t *testing.T) {
	defer RecoverTest(t)

	AssertEqual(try(func() { Assert(true, "aaaa") }), nil)
	Assert(false, "aaaa")
}

func try(fn func()) (err error) {
	defer RecoverErr(&err)

	fn()
	return nil
}

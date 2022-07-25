package assert

import (
	"testing"
)

func TestCheckNil(t *testing.T) {
	var a *int
	Assert(a == nil, "ok")
}

func try(fn func()) (err error) {
	fn()
	return nil
}

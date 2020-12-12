package xerror

import (
	"testing"
)

func TestCheckNil(t *testing.T) {
	defer RespDebug()

	var a *int
	CheckNil(a,"aaaa")
}

func TestCheck(t *testing.T) {
	defer RespDebug()

	Check(true,"aaaa")
	Check(false,"aaaa")
}

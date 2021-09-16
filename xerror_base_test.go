package xerror

import (
	"fmt"
	"testing"
)

func TestFmt(t *testing.T) {
	fmt.Println(Wrap(Fmt("hello %s","error")))
}
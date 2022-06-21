package syncx

import (
	"github.com/pubgo/funk"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	var is = assert.New(t)

	var cc = <-GoChan(func() Value[string] {
		return OK("ok")
	})
	funk.If(cc.IsErr(), func() {
	})

	is.Equal(cc.Get(), "ok")
}

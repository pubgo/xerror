package syncx

import (
	"testing"

	"github.com/pubgo/funk"
	"github.com/pubgo/funk/typex"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	var is = assert.New(t)

	var cc = <-GoChan(func() typex.Result[string] {
		return typex.OK("ok")
	})
	funk.If(cc.IsErr(), func() {
	})

	is.Equal(cc.Get(), "ok")
}

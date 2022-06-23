package syncx

import (
	"github.com/pubgo/funk"
	"github.com/pubgo/funk/typex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	var is = assert.New(t)

	var cc = <-GoChan(func() typex.Value[string] {
		return typex.OK("ok")
	})
	funk.If(cc.IsErr(), func() {
	})

	is.Equal(cc.Get(), "ok")
}

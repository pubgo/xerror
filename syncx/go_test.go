package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	var is = assert.New(t)

	var cc = <-GoChan(func() *Value[string] {
		return OK("ok")
	})

	is.Equal(cc.Val(), "ok")
}

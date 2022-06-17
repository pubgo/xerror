package typex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	var is = assert.New(t)
	var m Map[string]
	m.Set("a", "b")
	is.Equal(m.Get("a"), "b")
	t.Log(m.Map())
}

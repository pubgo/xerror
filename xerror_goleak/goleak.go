package xerror_goleak

import (
	"go.uber.org/goleak"

	"testing"
)

type Option = goleak.Option

var IgnoreTopFunction = goleak.IgnoreTopFunction
var IgnoreCurrent = goleak.IgnoreCurrent

func RespTest(t *testing.T, options ...Option) {
	if err := goleak.Find(options...); err != nil {
		t.Error(err)
	}
}

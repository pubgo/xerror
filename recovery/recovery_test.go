package recovery

import (
	"testing"

	"github.com/pubgo/funk/logx"
	"github.com/pubgo/funk/xerr"
)

func TestName(t *testing.T) {
	defer Recovery(func(err xerr.XErr) {
		err.DebugPrint()
	})

	logx.Info("test panic")
	hello()
}

func hello() {
	panic("hello")
}

package funk

import (
	"github.com/pubgo/funk/xerr"
	"log"
	"testing"
)

func TestName(t *testing.T) {
	defer Recovery(func(err xerr.XErr) {
		err.DebugPrint()
	})

	log.Println("test panic")
	hello()
}

func hello() {
	panic("hello")
}

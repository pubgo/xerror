package funk

import (
	"log"
	"testing"
)

func TestName(t *testing.T) {
	defer Recovery(func(err XErr) {
		err.DebugPrint()
	})

	log.Println("test panic")
	hello()
}

func hello() {
	panic("hello")
}

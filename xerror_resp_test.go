package xerror

import (
	"log"
	"testing"
)

func TestName(t *testing.T) {
	defer Resp(func(err XErr) {
		err.Debug()
	})

	log.Println("test panic")
	hello()
}

func hello() {
	panic("hello")
}

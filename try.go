package xerror

import (
	"fmt"
	"github.com/pubgo/xerror/internal/utils"
)

func TryCatch(fn func(), catch ...func(err error)) {
	Assert(fn == nil, "[fn] should not be nil")

	if len(catch) > 0 && catch[0] != nil {
		defer Resp(func(err XErr) { catch[0](err) })
	}

	fn()
}

func TryWith(err *error, fn func()) {
	Assert(fn == nil, "[fn] should not be nil")

	defer RespErr(err)
	fn()

	return
}

func TryThrow(fn func(), args ...interface{}) {
	Assert(fn == nil, "[fn] should not be nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err error
		handleRecover(&err, val)
		if isErrNil(err) {
			return
		}

		err1 := &xerror{Err: err}
		err1.Msg = fmt.Sprint(args...)
		err1.Caller = [2]string{utils.CallerWithFunc(fn)}
		panic(err1)
	}()

	fn()

	return
}

func Try(fn func()) (err error) {
	Assert(fn == nil, "[fn] should not be nil")

	defer RespErr(&err)

	fn()
	return
}

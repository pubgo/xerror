package xerror

import (
	"fmt"

	"github.com/pubgo/xerror/internal/utils"
)

func TryWith(err *error, fn func()) {
	if fn == nil {
		panic("[fn] should not be nil")
	}

	defer RespErr(err)
	fn()

	return
}

func TryCatch(fn func() (interface{}, error), catch func(err error)) interface{} {
	if fn == nil {
		panic("[fn] should not be nil")
	}

	if catch == nil {
		panic("[catch] should not be nil")
	}

	defer Resp(func(err XErr) { catch(err) })

	var val, err = fn()
	if err != nil {
		catch(err)
		return nil
	}
	return val
}

func TryVal(fn func() interface{}) interface{} {
	if fn == nil {
		panic("[fn] should not be nil")
	}

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
		err1.Caller = [2]string{
			utils.CallerWithFunc(fn),
			utils.CallerWithDepth(4),
		}
		panic(err1)
	}()

	return fn()
}

func TryThrow(fn func(), args ...interface{}) {
	if fn == nil {
		panic("[fn] should not be nil")
	}

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
		err1.Caller = [2]string{
			utils.CallerWithFunc(fn),
			utils.CallerWithDepth(4),
		}
		panic(err1)
	}()

	fn()

	return
}

func Try(fn func()) (err error) {
	if fn == nil {
		panic("[fn] should not be nil")
	}

	defer RespErr(&err)

	fn()
	return
}

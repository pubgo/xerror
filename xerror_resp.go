package xerror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/pubgo/xerror/internal/utils"
)

func RespErr(err *error) {
	val := recover()
	if val == nil {
		return
	}

	handleRecover(err, val)
}

func Raise(fns ...func(err XErr) error) {
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
	if len(fns) > 0 {
		err1.Caller = [2]string{utils.CallerWithFunc(fns[0])}
		panic(fns[0](err1))
	}

	panic(err1)
}

func RespRaise(fns ...func(err XErr) error) {
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
	if len(fns) > 0 {
		err1.Caller = [2]string{utils.CallerWithFunc(fns[0])}
		panic(fns[0](err1))
	}

	panic(err1)
}

func Resp(fn func(err XErr)) {
	Assert(fn == nil, "[fn] should not be nil")

	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	fn(&xerror{Err: err, Caller: [2]string{utils.CallerWithFunc(fn)}})
}

func RespExit(args ...interface{}) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }).stackString())
	printStack()
	os.Exit(1)
}

func RespDebug(args ...interface{}) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	p(handle(err, func(err *xerror) { err.Msg = fmt.Sprint(args...) }).stackString())
	printStack()
}

func RespTest(t *testing.T, debugs ...bool) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	var msg = handle(err).stackString()

	if len(debugs) > 0 {
		p(msg)
		return
	}

	t.Fatal(msg)
}

func RespHttp(w http.ResponseWriter, req *http.Request) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	w.WriteHeader(http.StatusBadRequest)

	var dt = PanicBytes(json.MarshalIndent(req.Header, "", "\t"))
	fmt.Fprintln(w, "request header")
	fmt.Fprintln(w, string(dt))
	fmt.Fprint(w, "\n\n\n\n")
	fmt.Fprintln(w, "error stack")
	fmt.Fprintln(w, handle(err).stackString())
	fmt.Fprint(w, "\n\n\n\n")
	fmt.Fprintln(w, "stack")
	buf := make([]byte, 1024*1024)
	fmt.Fprintln(w, string(buf[:runtime.Stack(buf, true)]))
}

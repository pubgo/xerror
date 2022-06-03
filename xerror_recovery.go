package xerror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/pubgo/xerror/internal/utils"
	"github.com/pubgo/xerror/xerror_core"
)

func RecoverErr(gErr *error) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	*gErr = &xerror{
		Err: err,
		Msg: err.Error(),
		Caller: []string{
			utils.CallerWithDepth(xerror_core.Conf.CallDepth + 1),
			utils.CallerWithDepth(xerror_core.Conf.CallDepth + 2),
		},
	}
}

func RecoverAndRaise(fns ...func(err XErr) XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	err1 := &xerror{
		Err: err,
		Msg: err.Error(),
	}
	if len(fns) > 0 {
		err1.Caller = []string{
			utils.CallerWithDepth(xerror_core.Conf.CallDepth + 1),
			utils.CallerWithDepth(xerror_core.Conf.CallDepth + 2),
		}
		panic(fns[0](err1))
	}
	panic(err1)
}

func Recovery(fn func(err XErr)) {
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

	fn(&xerror{
		Err: err,
		Msg: err.Error(),
		Caller: []string{
			utils.CallerWithDepth(xerror_core.Conf.CallDepth + 1),
			utils.CallerWithDepth(xerror_core.Conf.CallDepth + 2),
		}})
}

func RecoverAndExit() {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	var caller = []string{
		utils.CallerWithDepth(xerror_core.Conf.CallDepth + 2),
		utils.CallerWithDepth(xerror_core.Conf.CallDepth + 3),
	}

	p(handle(err, func(err *xerror) {
		err.Caller = append(err.Caller, caller...)
	}).debugString())
	printStack()
	os.Exit(1)
}

func RecoverTest(t *testing.T, debugs ...bool) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	var msg = handle(err).debugString()

	if len(debugs) > 0 {
		p(msg)
		return
	}

	t.Fatal(msg)
}

func RecoverHttp(w http.ResponseWriter, req *http.Request, fns ...func(err error)) {
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
	fmt.Fprintln(w, handle(err).debugString())
	fmt.Fprint(w, "\n\n\n\n")
	fmt.Fprintln(w, "stack")
	buf := make([]byte, 1024*1024)
	if len(fns) > 0 {
		fns[0](err)
	}

	fmt.Fprintln(w, string(buf[:runtime.Stack(buf, true)]))
}

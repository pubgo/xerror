package xerror

import (
	"os"
	"testing"

	"github.com/pubgo/xerror/internal/utils"
	"github.com/pubgo/xerror/xerror_core"
)

func RecoverErr(gErr *error, fns ...func(err XErr) XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	handleRecover(&err, val)
	if isErrNil(err) {
		return
	}

	err1 := handle(err)
	if len(fns) > 0 {
		*gErr = fns[0](err1)
		return
	}
	*gErr = err1
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

	err1 := handle(err)
	if len(fns) > 0 {
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

	fn(handle(err))
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

	err1 := handle(err)
	for i := 0; ; i++ {
		var cc = utils.CallerWithDepth(xerror_core.Conf.CallDepth + i)
		if cc == "" {
			break
		}
		err1.Caller = append(err1.Caller, cc)
	}

	p(err1.debugString())
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

	err1 := handle(err)
	for i := 0; ; i++ {
		var cc = utils.CallerWithDepth(xerror_core.Conf.CallDepth + i)
		if cc == "" {
			break
		}
		err1.Caller = append(err1.Caller, cc)
	}

	var msg = err1.debugString()

	if len(debugs) > 0 {
		p(msg)
		return
	}

	t.Fatal(msg)
}

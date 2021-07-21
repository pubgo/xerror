package xerror

import (
	"fmt"
	"reflect"

	"github.com/pubgo/xerror/internal/utils"
	"github.com/pubgo/xerror/xerror_core"
)

func Fmt(format string, a ...interface{}) *xerrorBase {
	x := &xerrorBase{}
	x.Code = fmt.Sprintf(format, a...)
	x.Caller = utils.CallerWithDepth(xerror_core.Conf.CallDepth+1)
	return x
}

func New(code string, ms ...string) *xerrorBase {
	var msg string
	if len(ms) > 0 {
		msg = ms[0]
	}

	xw := &xerrorBase{}
	xw.Code = code
	xw.Msg = msg
	xw.Caller = utils.CallerWithDepth(xerror_core.Conf.CallDepth)

	return xw
}

type xerrorBase struct {
	Code   string `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Caller string `json:"caller,omitempty"`
}

func (t *xerrorBase) Error() string { return fmt.Sprintf("[%s] %s", t.Code, t.Msg) }
func (t *xerrorBase) Is(err error) bool {
	if t == nil || err == nil {
		return false
	}

	switch err := err.(type) {
	case *xerrorBase:
		return err.Code == t.Code
	default:
		return false
	}
}
func (t *xerrorBase) As(target interface{}) bool {
	t1 := reflect.Indirect(reflect.ValueOf(target)).Interface()
	if err, ok := t1.(*xerrorBase); ok {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

func (t *xerrorBase) New(msg string) error {
	x := &xerrorBase{Code: t.Code}
	x.Msg = msg
	x.Caller = utils.CallerWithDepth(xerror_core.Conf.CallDepth + 1)
	return x
}

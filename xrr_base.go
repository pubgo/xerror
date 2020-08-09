package xerror

import (
	"github.com/pubgo/xerror/xerror_util"
	"strings"

	"github.com/pubgo/xerror/internal/wrapper"
	"github.com/pubgo/xerror/xerror_core"
)

type xerrorBase struct {
	Code   string `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Caller string `json:"caller,omitempty"`
}

func (t *xerrorBase) Error() string {
	return t.Code
}

func (t *xerrorBase) As(err interface{}) bool {
	if t == nil || err == nil {
		return false
	}

	switch e := err.(type) {
	case **xerrorBase:
		return strings.HasPrefix(t.Code, (*e).Code)
	case *xerrorBase:
		return strings.HasPrefix(t.Code, e.Code)
	case *error:
		return strings.HasPrefix(t.Code, (*e).Error())
	case error:
		return strings.HasPrefix(t.Code, e.Error())
	case *string:
		return strings.HasPrefix(t.Code, *e)
	case string:
		return strings.HasPrefix(t.Code, e)
	default:
		return false
	}
}

func (t *xerrorBase) FamilyAs(err interface{}) bool {
	return t.As(err)
}

func (t *xerrorBase) New(code string, ms ...string) error {
	var msg string
	if len(ms) > 0 {
		msg = ms[0]
	}

	code = t.Code + xerror_core.Delimiter + code
	xw := &xerrorBase{}
	xw.Code = code
	xw.Msg = msg
	xw.Caller = xerror_util.CallerWithDepth(wrapper.CallDepth())

	return xw
}

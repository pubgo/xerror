package funk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pubgo/funk/funkonf"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/internal/utils"
)

const callStackDepth = 2

func NewErr(format string, a ...interface{}) XErr {
	x := &xerror{}
	x.Msg = fmt.Sprintf(format, a...)
	x.Caller = []string{utils.CallerWithDepth(callStackDepth + 1)}
	return x
}

type xerror struct {
	Err    error    `json:"cause,omitempty"`
	Msg    string   `json:"msg,omitempty"`
	Detail string   `json:"detail,omitempty"`
	Caller []string `json:"caller,omitempty"`
}

func (t *xerror) xErr()          {}
func (t *xerror) String() string { return t.Stack() }
func (t *xerror) DebugPrint() {
	if !funkonf.Conf.Debug {
		return
	}

	p(handle(Wrap(t)).debugString())
}

func (t *xerror) Unwrap() error { return t.Err }
func (t *xerror) Cause() error  { return t.Err }
func (t *xerror) Wrap(args ...interface{}) XErr {
	return handle(t, func(err *xerror) { err.Detail = fmt.Sprint(args...) })
}

func (t *xerror) WrapF(msg string, args ...interface{}) XErr {
	return handle(t, func(err *xerror) { err.Detail = fmt.Sprintf(msg, args...) })
}

func (t *xerror) _p(buf *strings.Builder, xrr *xerror) {
	if xrr == nil {
		return
	}

	buf.WriteString("========================================================================================================================\n")
	if xrr.Msg != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Green.P("Msg"), xrr.Msg))
	}

	if xrr.Detail != "" {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Green.P("Detail"), xrr.Detail))
	}

	for i := range xrr.Caller {
		if strings.Contains(xrr.Caller[i], "/src/runtime/") {
			continue
		}
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Yellow.P("Caller"), xrr.Caller[i]))
	}

	t._p(buf, trans(xrr.Err))
}

func (t *xerror) debugString() string {
	if t == nil || t.Err == nil {
		return ""
	}

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString("\n")
	t._p(buf, t)
	buf.WriteString("========================================================================================================================\n\n")
	return buf.String()
}

func (t *xerror) Is(err error) bool {
	if t == nil || t.Err == nil || err == nil {
		return false
	}

	switch _err := err.(type) {
	case *xerror:
		return _err == t || _err.Err == t.Err
	case error:
		return t.Err == _err
	default:
		return false
	}
}

func (t *xerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('#') {
			type errors xerror
			_, _ = fmt.Fprintf(s, "%#v", (*errors)(t))
			return
		}

		if s.Flag('+') {
			_, _ = fmt.Fprint(s, t.Stack())
			return
		}

		_, _ = fmt.Fprint(s, t.Stack())
	case 's', 'q':
		_, _ = fmt.Fprint(s, t.Msg+": \n\t"+t.Error()+"\n\t"+t.Caller[0]+"\n\t"+t.Caller[1])
	default:
		_, _ = fmt.Fprint(s, t.Msg)
	}
}

func (t *xerror) Stack() string {
	if t == nil || t.Err == nil {
		return ""
	}

	dt, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(dt)
}

func (t *xerror) As(target interface{}) bool {
	if t == nil || target == nil {
		return false
	}

	var v = reflect.ValueOf(target)
	t1 := reflect.Indirect(v).Interface()
	if err, ok := t1.(*xerror); ok {
		v.Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

// Error
func (t *xerror) Error() string {
	if t == nil || isErrNil(t.Err) {
		return ""
	}

	return t.Err.Error()
}

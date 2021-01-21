package xerror

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pubgo/xerror/internal/color"
	"github.com/pubgo/xerror/xerror_abc"
)

type XErr = xerror_abc.XErr
type xerror struct {
	Cause1 error     `json:"cause,omitempty"`
	Msg    string    `json:"msg,omitempty"`
	Caller [2]string `json:"caller,omitempty"`
}

func (t *xerror) Wrap(args ...interface{}) error              { return With().Wrap(t, args...) }
func (t *xerror) WrapF(msg string, args ...interface{}) error { return With().WrapF(t, msg, args...) }
func (t *xerror) Unwrap() error                               { return t.Cause() }
func (t *xerror) Cause() error {
	if t == nil {
		return nil
	}

	return t.Cause1
}

func (t *xerror) _p(buf *strings.Builder, xrr *xerror) {
	buf.WriteString("========================================================================================================================\n")
	if xrr.Cause1 != nil {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Red.P("Err"), xrr.Cause1.Error()))
	}
	if strings.TrimSpace(xrr.Msg) != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Green.P("Msg"), xrr.Msg))
	}

	for i := range xrr.Caller {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Yellow.P("Caller"), xrr.Caller[i]))
	}

	if errs := trans(xrr.Cause1); errs != nil {
		for i := range errs {
			t._p(buf, errs[i])
		}
	}
}

func (t *xerror) p() string {
	if t == nil || t.Cause1 == nil {
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
	if t == nil || t.Cause1 == nil || err == nil {
		return false
	}

	switch err := err.(type) {
	case *xerrorBase:
		return err == t.Cause1
	case *xerror:
		return err == t || err.Cause1 == t.Cause1
	case error:
		return t.Cause1 == err
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
			_, _ = fmt.Fprint(s, t.Stack(true))
			return
		}

		_, _ = fmt.Fprint(s, t.Stack())
	case 's', 'q':
		_, _ = fmt.Fprint(s, t.Msg+": \n\t"+t.Error()+"\n\t"+t.Caller[0]+"\n\t"+t.Caller[1])
	default:
		_, _ = fmt.Fprint(s, t.Msg)
	}
}

func (t *xerror) Stack(indent ...bool) string {
	if t == nil || t.Cause1 == nil || t.Cause1 == ErrDone {
		return ""
	}

	if len(indent) > 0 {
		dt, err := json.MarshalIndent(t, "", "\t")
		if err != nil {
			log.Fatalln(err)
		}
		return string(dt)
	}
	dt, err := json.Marshal(t)
	if err != nil {
		log.Fatalln(err)
	}
	return string(dt)
}

// Error
func (t *xerror) Error() string {
	if t == nil || isErrNil(t.Cause1) {
		return ""
	}

	return t.Cause1.Error()
}

func (t *xerror) String() string {
	return t.Stack()
}

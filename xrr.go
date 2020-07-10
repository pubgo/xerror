package xerror

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type xerrorBase struct {
	Code  string `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Caller string `json:"caller,omitempty"`
}

func (t *xerrorBase) Error() string {
	return t.Code
}

func (t *xerrorBase) New(code string, ms ...string) error {
	var msg string
	if len(ms) > 0 {
		msg = ms[0]
	}

	code = t.Code + "->" + code
	xw := &xerrorBase{}
	xw.Code = code
	xw.Msg = msg
	xw.Caller = callerWithDepth(callDepth)

	return xw
}

type xerror struct {
	Cause1 error  `json:"next,omitempty"`
	Code1  string `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Caller string `json:"caller,omitempty"`
}

func (t *xerror) Code() string {
	return t.Code1
}

func (t *xerror) Unwrap() error {
	if t == nil {
		return nil
	}

	return t.Cause1
}

func (t *xerror) Cause() error {
	if t == nil {
		return nil
	}

	return t.Cause1
}

func (t *xerror) p() string {
	if t == nil || t.Cause1 == nil {
		return ""
	}

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString("\n")
	xrr := t
	for xrr != nil {
		buf.WriteString("========================================================================================================================\n")
		if xrr.Cause1 != nil {
			buf.WriteString(fmt.Sprintf("   %s]: %s\n", colorize("Err", colorRed), xrr.Cause1))
		}
		if xrr.Msg != "" {
			buf.WriteString(fmt.Sprintf("   %s]: %s\n", colorize("Msg", colorGreen), xrr.Msg))
		}
		if xrr.Code1 != "" {
			buf.WriteString(fmt.Sprintf("  %s]: %s\n", colorize("Code", colorGreen), xrr.Code1))
		}
		buf.WriteString(fmt.Sprintf("%s]: %s\n", colorize("Caller", colorYellow), xrr.Caller))
		xrr = trans(xrr.Cause1)
	}
	buf.WriteString("========================================================================================================================\n\n")
	return buf.String()
}

func (t *xerror) Is(err error) bool {
	if t == nil || t.Cause1 == nil || err == nil {
		return false
	}

	switch err := err.(type) {
	case *xerror:
		return err == t || err.Cause1 == t.Cause1 || err.Code1 == t.Code1
	case error:
		return t.Cause1 == err
	default:
		return false
	}
}

func (t *xerror) As(err interface{}) bool {
	if t == nil || t.Cause1 == nil || err == nil {
		return false
	}

	switch e := err.(type) {
	case *xerror:
		return strings.HasPrefix(t.Code1, e.Code1)
	case error:
		return strings.HasPrefix(t.Code1, e.Error())
	case string:
		return strings.HasPrefix(t.Code1, e)
	default:
		return false
	}
}

func (t *xerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", t.Cause())
			io.WriteString(s, t.Msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, t.Error())
	}
}

func (t *xerror) Stack() string {
	if t == nil || t.Cause1 == nil || t.Cause1 == ErrDone {
		return ""
	}
	dt, _ := json.Marshal(t)
	return string(dt)
}

// Error
func (t *xerror) Error() string {
	if t == nil || t.Cause1 == nil || t.Cause1 == ErrDone {
		return ""
	}
	return t.Cause1.Error()
}

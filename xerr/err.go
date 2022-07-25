package xerr

import (
	"errors"
	"fmt"
)

type Err struct {
	Err    error  `json:"err"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
}

func (e Err) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}

	return errors.New(e.String())
}

func (e Err) String() string {
	return fmt.Sprintf("msg=%s detail=%s", e.Msg, e.Detail)
}

func (e Err) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.String()
}

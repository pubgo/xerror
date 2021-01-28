package xerror

import (
	"strings"
)

type combine []*xerror

func (errs combine) String() string {
	var b strings.Builder
	defer b.Reset()
	for i := range errs {
		b.WriteString(errs[i].String())
		b.WriteString("\n")
	}
	return b.String()
}

func (errs combine) Error() string {
	var b strings.Builder
	defer b.Reset()
	for i := range errs {
		b.WriteString(errs[i].Error())
		b.WriteString("\n")
	}
	return b.String()
}

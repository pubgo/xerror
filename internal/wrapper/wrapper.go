package wrapper

import (
	"runtime/debug"

	"github.com/pubgo/xerror/xerror_core"
)

func IsCaller() bool {
	return xerror_core.IsCaller
}
func CallDepth() int {
	return xerror_core.CallDepth
}

func PrintStack() {
	if xerror_core.PrintStack {
		debug.PrintStack()
	}
}

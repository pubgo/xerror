package xerror_core

import (
	"os"
	"strconv"

	"github.com/pubgo/xerror/xerror_envs"
)

const Delimiter = "||"

var IsCaller bool
var CallDepth = 3
var PrintStack bool

func init() {
	IsCaller = true
	PrintStack = true

	if env := os.Getenv(xerror_envs.IsCaller); env != "" {
		IsCaller, _ = strconv.ParseBool(env)
	}

	if env := os.Getenv(xerror_envs.PrintStack); env != "" {
		PrintStack, _ = strconv.ParseBool(env)
	}
}

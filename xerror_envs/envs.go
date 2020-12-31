package xerror_envs

import (
	"expvar"
	"os"
	"strconv"
	"strings"

	"github.com/pubgo/xerror/internal/envs"
)

var prefix = "XERROR_"

const (
	IsCaller   = "XERROR_CALLER"
	PrintStack = "XERROR_STACK"
	CallDepth  = "XERROR_CALLDEPTH"
	Delimiter  = "XERROR_DELIMITER"
)

func IsCallerVal() bool    { return envs.IsCaller }
func PrintStackVal() bool  { return envs.PrintStack }
func CallDepthVal() int    { return envs.CallDepth }
func DelimiterVal() string { return envs.Delimiter }

func init() {
	if env := os.Getenv(IsCaller); env != "" {
		envs.IsCaller, _ = strconv.ParseBool(env)
	}

	if env := os.Getenv(PrintStack); env != "" {
		envs.PrintStack, _ = strconv.ParseBool(env)
	}

	if env := os.Getenv(CallDepth); env != "" {
		envs.CallDepth, _ = strconv.Atoi(env)
	}

	if env := os.Getenv(Delimiter); env != "" {
		envs.Delimiter = env
	}

	expvar.Publish("xerror_envs", expvar.Func(func() interface{} {
		var data []string
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, prefix) {
				data = append(data, env)
			}

		}
		return data
	}))
}

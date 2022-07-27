package logx

import (
	"github.com/go-logr/logr"
)

func WithCallDepth(depth int) logr.Logger {
	return logT.WithCallDepth(depth)
}

func WithName(name string) logr.Logger {
	return logT.WithName(name)
}

func V(level int) logr.Logger {
	return logT.V(level)
}

func WithValues(keysAndValues ...interface{}) logr.Logger {
	return logT.WithValues(keysAndValues...)
}

func IfEnabled(level int, fn func(log logr.Logger)) {
	var log = V(level)
	if log.Enabled() {
		fn(log)
	}
}

func Enabled() bool {
	return logT.Enabled()
}

func Info(msg string, keysAndValues ...interface{}) {
	logT.WithCallDepth(1).Info(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	if err == nil {
		return
	}

	logT.WithCallDepth(1).Error(err, msg, keysAndValues...)
}

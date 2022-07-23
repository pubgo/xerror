package logx

import (
	"github.com/go-logr/logr"
)

func Named(name string) logr.Logger {
	return logT.WithName(name)
}

func V(level int) logr.Logger {
	return logT.V(level)
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

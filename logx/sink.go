package logx

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	logkit "github.com/go-kit/log"
	"github.com/go-logr/logr"
)

var _ logr.LogSink = (*sink)(nil)
var _ logr.CallDepthLogSink = (*sink)(nil)

type sink struct {
	level     int
	callDepth int
	prefix    string
	values    []interface{}
	log       *logr.Logger
}

// Enabled reports whether this Logger is enabled with respect to the current global log level.
func (s *sink) Enabled(level int) bool {
	if level > int(atomic.LoadInt32(&gv)) {
		return false
	}

	if s.prefix == "" {
		return true
	}

	return true
}

func (s sink) WithCallDepth(depth int) logr.LogSink {
	s.callDepth += depth
	return &s
}

func (s *sink) Init(info logr.RuntimeInfo) {
	s.callDepth += info.CallDepth
}

func (s *sink) Info(level int, msg string, keysAndValues ...interface{}) {
	if !s.Enabled(level) {
		return
	}

	if defaultLog == nil {
		s.log.WithCallDepth(s.callDepth).WithName(s.prefix).WithValues(s.values...).GetSink().Info(level, msg, keysAndValues...)
		return
	}

	keysAndValues = append(keysAndValues, s.values...)
	keysAndValues = append(keysAndValues, "caller", logkit.Caller(s.callDepth+DefaultCallerSkip)())
	keysAndValues = append(keysAndValues, "logger", s.prefix)
	keysAndValues = append(keysAndValues, "level", "info")
	keysAndValues = append(keysAndValues, "msg", msg)
	keysAndValues = append(keysAndValues, "ts", time.Now().UTC().Format(TimestampFormat))
	if err := defaultLog.Log(keysAndValues...); err != nil {
		panic(err)
	}
}

func (s *sink) Error(err error, msg string, keysAndValues ...interface{}) {
	if err == nil {
		return
	}

	if defaultLog == nil {
		keysAndValues = append(keysAndValues, "error_msg", fmt.Sprintf("%#v", err))
		keysAndValues = append(keysAndValues, "stacktrace", stringify(debug.Stack()))
		s.log.WithCallDepth(s.callDepth).WithName(s.prefix).WithValues(s.values...).GetSink().Error(err, msg, keysAndValues...)
		return
	}

	keysAndValues = append(keysAndValues, s.values...)
	keysAndValues = append(keysAndValues, "caller", logkit.Caller(s.callDepth+DefaultCallerSkip)())
	keysAndValues = append(keysAndValues, "logger", s.prefix)
	keysAndValues = append(keysAndValues, "level", "error")
	keysAndValues = append(keysAndValues, "msg", msg)
	keysAndValues = append(keysAndValues, "error", err.Error())
	keysAndValues = append(keysAndValues, "error_msg", fmt.Sprintf("%#v", err))
	keysAndValues = append(keysAndValues, "stacktrace", stringify(debug.Stack()))
	keysAndValues = append(keysAndValues, "ts", time.Now().UTC().Format(TimestampFormat))
	if err := defaultLog.Log(keysAndValues...); err != nil {
		panic(err)
	}
}

func (s sink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	s.values = append(s.values, keysAndValues...)
	return &s
}

func (s sink) WithName(name string) logr.LogSink {
	if len(s.prefix) > 0 {
		s.prefix = s.prefix + "."
	}
	s.prefix += name
	return &s
}

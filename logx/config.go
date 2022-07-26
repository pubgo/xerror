package logx

import (
	"sync/atomic"

	logkit "github.com/go-kit/log"
	"github.com/go-logr/logr"
	"github.com/iand/logfmtr"
)

var defaultLog logkit.Logger
var logT = logr.Discard()
var TimestampFormat = "2006-01-02T15:04:05.000000000Z07:00"
var NameDelim = "."
var DefaultCallerSkip = 2

// The global verbosity level.
var gv int32 = 2

func init() {
	opts := logfmtr.DefaultOptions()
	opts.Humanize = true
	opts.Colorize = true
	opts.CallerSkip = DefaultCallerSkip
	opts.AddCaller = true

	var log = logfmtr.NewWithOptions(opts)
	logT = logr.New(&sink{log: &log})
}

func SetLog(log logkit.Logger) {
	defaultLog = log
}

// SetVerbosity sets the global log level. Only loggers with a V level less than
// or equal to this value will be enabled.
func SetVerbosity(v int) int {
	old := atomic.SwapInt32(&gv, int32(v))
	return int(old)
}

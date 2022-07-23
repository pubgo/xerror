package logx

import (
	"sync/atomic"

	"github.com/go-logr/logr"
	"github.com/iand/logfmtr"
)

var defaultLog = logr.Discard()
var changeNum int32 = 1
var logT = logr.New(&sink{})

func init() {
	opts := logfmtr.DefaultOptions()
	opts.Humanize = true
	opts.Colorize = true
	opts.CallerSkip = 2
	opts.AddCaller = true
	logfmtr.UseOptions(opts)
	defaultLog = logfmtr.New()
}

func SetLog(log logr.Logger) {
	defaultLog = log
	atomic.AddInt32(&changeNum, 1)
}

func SetVerbosity(v int) {
	logfmtr.SetVerbosity(v)
}

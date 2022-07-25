package logx

import (
	"github.com/go-logr/logr"
	"github.com/iand/logfmtr"
)

var defaultLog = logr.Discard()
var logT = logr.New(&sink{})

func init() {
	SetVerbosity(2)

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
}

func SetVerbosity(v int) {
	logfmtr.SetVerbosity(v)
}

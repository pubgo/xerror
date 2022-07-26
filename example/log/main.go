package main

import (
	"fmt"
	"os"

	logkit "github.com/go-kit/log"
	"github.com/go-logr/logr"

	"github.com/pubgo/funk/logx"
)

var dd = logx.WithName("dd")

func main() {
	demo(dd)
	demo(logx.V(1).WithName("abc"))
	logx.Info("test")
	demo(logx.WithName("demo"))

	var ll = logkit.NewJSONLogger(os.Stdout)
	logx.SetLog(ll)
	demo(logx.WithName("logkit"))
	demo(dd)
}

func demo(base logr.Logger) {
	l := base.WithName("MyName").WithName("dd").WithValues("user", "you")
	l.Info("hello", "val1", 1, "val2", map[string]int{"k": 1})
	l.V(1).Info("you should see this")
	l.V(1).V(1).Info("you should NOT see this")
	l.Error(nil, "uh oh", "trouble", true, "reasons", []float64{0.1, 0.11, 3.14})
	l.Error(fmt.Errorf("an error occurred"), "goodbye", "code", -1)
}

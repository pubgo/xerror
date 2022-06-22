package main

import (
	"github.com/pubgo/funk"
)

// 单个pkg的error处理

var err1 = &funk.Err{Msg: "业务错误处理", Detail: "详细信息"}

func Hello() {
	defer funk.RecoverAndRaise(func(err funk.XErr) funk.XErr {
		return err.Wrap("Hello wrap")
	})

	var err2 = funk.WrapF(err1, "处理 wrap")
	funk.MustF(err2, "处理 panic")
	return
}

func CallHello() (gErr error) {
	defer funk.Recovery(func(err funk.XErr) {
		gErr = err.WrapF("CallHello wrap")
	})

	Hello()

	return
}

func main() {
	defer funk.RecoverAndExit()

	funk.Must(CallHello())
}

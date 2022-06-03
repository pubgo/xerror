package main

import (
	"github.com/pubgo/xerror"
)

// 单个pkg的error处理

var err1 = &xerror.Err{Msg: "业务错误处理", Detail: "详细信息"}

func Hello() (err error) {
	defer xerror.RecoverErr(&err)

	xerror.Panic(err1, "处理 业务错误处理 失败")
	return
}

func CallHello() (gErr error) {
	defer xerror.Recovery(func(err xerror.XErr) {
		gErr = err.WrapF("CallHello wrap")
	})

	xerror.Panic(Hello())

	return
}

func main() {
	defer xerror.RecoverAndExit()

	xerror.Panic(CallHello())
}

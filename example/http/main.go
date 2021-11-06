package main

import (
	"fmt"
	"net/http"

	"github.com/pubgo/xerror"
)

func main() {
	xerror.Exit(http.ListenAndServe(":8088", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 处理http handler的panic, 并附带丰富的堆栈信息
		defer xerror.RespHttp(writer, request, func(err error) {
			fmt.Println(err)
		})
		xerror.Panic(fmt.Errorf("panic"))
	})))
}

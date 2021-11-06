package main

import (
	"fmt"
	"net/http"

	"github.com/pubgo/xerror"
)

func main() {
	xerror.Exit(http.ListenAndServe(":8088", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer xerror.RespHttp(writer, request, func(err error) {
			fmt.Println(err)
		})
		xerror.Panic(fmt.Errorf("panic"))
	})))
}

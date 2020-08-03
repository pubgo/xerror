# xerror

go error 简单实现


## 性能分析
```sh
go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out ./...
go tool pprof -http=":8081" profile.out
go tool pprof -http=":8081" memprofile.out
```

```md
goos: darwin
goarch: amd64
pkg: github.com/pubgo/xerror
BenchmarkPanic-8                	 1353934	       883 ns/op	     128 B/op	       2 allocs/op
BenchmarkPanicWithOutCaller-8   	 3861938	       309 ns/op	      48 B/op	       1 allocs/op
BenchmarkNoPanic-8              	201641330	         5.87 ns/op	       0 B/op	       0 allocs/op
PASS
ok      github.com/pubgo/xerror 4.363s
```

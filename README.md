# xerror

go error 简单实现


## 性能分析
```sh
go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out
go tool pprof -http=":8081" profile.out
go tool pprof -http=":8081" memprofile.out
```

```md
goos: darwin
goarch: amd64
pkg: github.com/pubgo/xerror
BenchmarkPanic-8          403794              2991 ns/op             382 B/op          7 allocs/op
BenchmarkNoPanic-8      188605891             6.46 ns/op             0 B/op            0 allocs/op
PASS
ok      github.com/pubgo/xerror 4.363s
```

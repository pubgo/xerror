
.PHONY: test
test:
	@go test -timeout=1s -race -cover -coverprofile=out.out ./...
	@go tool cover -func=out.out

.PHONY: test_html
test_html:
	@go test -timeout=1s -race -cover -coverprofile=out.out ./...
	@go tool cover -html=out.out

.PHONY: test_bench
test_bench:
	@go test -bench=. -benchmem ./

.PHONY: rm_test
rm_test:
	@rm -f *.out
	@rm -f *.test

.PHONY: test_profile
test_profile:
	@go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out example_test.go
	@go tool pprof -http=":8081" profile.out

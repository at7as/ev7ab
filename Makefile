run:
	@go run ./cmd/app

run-test:
	@go run ./cmd/app ./test/config.json

build-app:
	@go build -o ./build/ ./cmd/app

bench:
	@go test -bench=. -benchmem ./cmd/bench/bench_test.go

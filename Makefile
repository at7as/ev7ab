run:
	@go run ./cmd/app --config=./test/example_simple/app.config.json

run-test:
	@go run ./cmd/app --config=./test/example_simple/app.config.json

build-app:
	@go build -o ./build/ ./cmd/app

bench:
	@go test -bench=. -benchmem ./cmd/bench/bench_test.go

bench-pprof:
	@go test -bench=. -benchmem ./cmd/bench/bench_test.go -cpuprofile cpu.out

pprof:
	@go tool pprof -http 127.0.0.1:8080 cpu.out

try:
	@go run ./cmd/try

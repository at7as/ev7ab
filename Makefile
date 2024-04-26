run:
	@go run ./cmd/app

run-test:
	@go run ./cmd/app ./test/config.json

build-app:
	@go build -o ./build/ ./cmd/app

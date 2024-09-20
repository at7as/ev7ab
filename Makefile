run:
	@go run ./cmd/app --config=./test/example_digits/app.config.json

build-app:
	@go build -o ./build/ ./cmd/app

try:
	@go run ./cmd/try

app:
	@./build/app --config=./test/example_digits/app.config.json
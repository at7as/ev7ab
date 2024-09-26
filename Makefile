simple-app:
	@go run ./examples/simple/main.go -config=./examples/simple/app.config.json

simple-try:
	@go run ./examples/simple/main.go -try

bezier-app:
	@go run ./examples/bezier/main.go -config=./examples/bezier/app.config.json

bezier-try:
	@go run ./examples/bezier/main.go -try

track-app:
	@go run ./examples/track/main.go -config=./examples/track/app.config.json

track-try:
	@go run ./examples/track/main.go -try

tictactoe-app:
	@go run ./examples/tictactoe/main.go -config=./examples/tictactoe/app.config.json

tictactoe-try:
	@go run ./examples/tictactoe/main.go -try

digits-app:
	@go run ./examples/digits/main.go -config=./examples/digits/app.config.json

digits-try:
	@go run ./examples/digits/main.go -try

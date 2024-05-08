run:
	@go run ./cmd/main/main.go

build:
	@go build -o ./bin/fetch-websocket ./cmd/main/main.go
build:
	@go build ./cmd/.
	@echo services were built!

run: build
	@go run ./cmd/main.go
	@echo app is running!
test:
	@echo "Testing..."
	go test -v ./...
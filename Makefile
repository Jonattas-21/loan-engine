build:
	@go build ./cmd/.
	@echo services were built!
	@golangci-lint run
	@echo linter passed!


run: build
	@go run ./cmd/main.go
	@echo app is running!
test:
	@echo "Testing..."
	go test -v ./...

lint:
	@echo "Linting..."
	golangci-lint run
	@echo linter passed!
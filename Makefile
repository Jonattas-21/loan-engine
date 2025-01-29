build:
	@go build ./cmd/.
	@echo services were built!
	@golangci-lint run
	@echo linter passed!
	@echo It's running...

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

swagger:
	@echo "Generating swagger..."
	swag init -g cmd/main.go --parseDependency --parseInternal
	@echo swagger generated!
build:
	@go build ./cmd/.
	@echo services were built!
	@golangci-lint run
	@echo linter passed!


run: build
	@go run ./cmd/main.go
	@echo It's running...
	
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
build:
	@go build ./cmd/.
	@echo services were built!
	@golangci-lint run
	@echo linter passed!

run: build
	@echo It's running...
	@go run ./cmd/main.go
	
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

# Access redis:
# docker exec -it local-redis redis-cli
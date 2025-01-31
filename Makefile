build:
	@go build ./cmd/.
	@echo services were built!
	@golangci-lint run
	@echo linter passed!

run: build
	@echo It's running...
	@go run ./cmd/main.go
	
lint:
	@echo "Linting..."
	golangci-lint run
	@echo linter passed!

swagger:
	@echo "Generating swagger..."
	swag init -g cmd/main.go --parseDependency --parseInternal
	@echo swagger generated!

test:
	@echo "Running tests..."
	@go test -v X:\Source\loan-engine\internal\usecases
	@echo tests passed!

test-create-profile:
	@go test X:\Source\loan-engine\internal\usecases -coverprofile=X:\Source\loan-engine\tests\test_profile
	@echo profile created!

test-coverage:
	@go test -cover X:\Source\loan-engine\internal\usecases
	@echo coverage passed!

test-html:
	@go tool cover -html=X:\Source\loan-engine\tests\test_profile
	@echo html passed!


# Access redis:
# docker exec -it local-redis redis-cli
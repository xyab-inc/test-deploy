.PHONY: test test-coverage build clean

# Default target
all: test build

# Run tests with coverage
test:
	go test -v -coverprofile=coverage.out ./...

# Show test coverage in browser
test-coverage: test
	go tool cover -html=coverage.out

# Show test coverage in terminal
test-coverage-func: test
	go tool cover -func=coverage.out

# Build the binary
build:
	go build -o bin/action

# Clean build artifacts
clean:
	rm -f coverage.out
	rm -rf bin/

# Install dependencies
deps:
	go mod download

# Run format and lint
fmt:
	go fmt ./...

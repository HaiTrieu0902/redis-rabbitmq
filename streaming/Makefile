.PHONY: build run clean test deps

# Build the application
build:
	go build -o bin/streaming main.go

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/streaming main.go

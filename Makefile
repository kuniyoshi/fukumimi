.PHONY: build run clean test fmt vet staticcheck lint check install build-all

# Binary name
BINARY_NAME=fukumimi

# Build the binary
build:
	go build -o $(BINARY_NAME) main.go

# Run the application
run:
	go run main.go

# Run with login command
run-login:
	go run main.go login

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run staticcheck (requires staticcheck to be installed)
# Install: go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck:
	staticcheck ./...

# Run linter (requires golangci-lint)
# Install: https://golangci-lint.run/usage/install/
lint:
	golangci-lint run

# Run all checks (fmt, vet, staticcheck)
check: fmt vet staticcheck

# Install the binary to GOPATH/bin
install:
	go install

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe main.go
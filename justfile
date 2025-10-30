# Default recipe to display available commands
default:
    @just --list

# Run the server
serve:
    go run main.go

# Run tests
test:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -v -cover -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

# Build the binary
build:
    go build -o testdummy -ldflags="-X main.versionString=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" .

# Build for multiple platforms
build-all:
    GOOS=linux GOARCH=amd64 go build -o testdummy-linux-amd64 -ldflags="-X main.versionString=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" .
    GOOS=darwin GOARCH=amd64 go build -o testdummy-darwin-amd64 -ldflags="-X main.versionString=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" .
    GOOS=darwin GOARCH=arm64 go build -o testdummy-darwin-arm64 -ldflags="-X main.versionString=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" .
    GOOS=windows GOARCH=amd64 go build -o testdummy-windows-amd64.exe -ldflags="-X main.versionString=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" .

# Clean build artifacts
clean:
    rm -f testdummy testdummy-* coverage.out coverage.html

# Run go mod tidy
tidy:
    go mod tidy

# Run linters (requires staticcheck)
lint:
    go vet ./...
    staticcheck ./...

# Format code
fmt:
    go fmt ./...

# Run the server with custom bind address
serve-custom ADDR:
    TESTDUMMY_BIND_ADDRESS={{ADDR}} go run main.go

# Run the server with debug logging
serve-debug:
    TESTDUMMY_ENABLE_REQUEST_LOGGING=true go run main.go

# Build and run
run: build
    ./testdummy

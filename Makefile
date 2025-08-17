OUT := gps-microservice
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ./...)  # all packages in the module

# Target architecture for the Pi
GOOS := linux
GOARCH := arm64  # Change to "arm" for 32-bit OS

# Output binary name with version and architecture
OUT_NAME := ${OUT}-$(VERSION)-${GOARCH}

all: build

# Build Linux binary for Pi (cross-compile from Windows)
build:
	@echo "Cross-compiling for Linux (${GOARCH})"
	@cmd /C "set GOOS=${GOOS}&& set GOARCH=${GOARCH}&& set CGO_ENABLED=0&& go build -v -o ${OUT_NAME} -ldflags=\"-X main.version=${VERSION}\" ./"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ${PKG_LIST}

# Vet code
vet:
	@echo "Running go vet..."
	@go vet ${PKG_LIST}

# Lint code
lint:
	@echo "Running golint..."
	@golint ${PKG_LIST} || exit 1

# Clean output binaries
clean:
	@cmd /C "if exist ${OUT_NAME} del ${OUT_NAME}"

# Run all checks
check: vet lint test
	@echo "All checks passed!"

.PHONY: all build test vet lint clean check

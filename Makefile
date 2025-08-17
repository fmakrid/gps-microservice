OUT := gps-microservice
PKG := github.com/fmakrid/gps-microservice
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/...)
GO_FILES := $(shell find . -name "*.go")

# Output binary name with version
OUT_NAME := gps-microservice-$(VERSION)

all: build

# Build for Ubuntu Linux (Linux amd64)
build:
	@echo "Building for Linux (Ubuntu)"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o ${OUT_NAME} -ldflags="-X main.version=${VERSION}" ${PKG}

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@golint ${PKG}/... || exit 1

clean:
	@rm -f ${OUT_NAME} ${OUT_NAME}-*

.PHONY: all build test vet lint clean

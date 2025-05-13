OUT := gps-microservice
PKG := github.com/fmakrid/gps-microservice
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | findstr /V "\\vendor")
GO_FILES := $(shell dir /s /b *.go | findstr /V "\\vendor")

all: run

server:
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build -v -o ${OUT} -ldflags='-X main.version=${VERSION}' ${PKG}"

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

# Lint target: run golint on all Go files in the package
lint:
	@golint ${PKG}/... || exit 1  # Use golint on the whole package

static: vet lint
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build -v -o ${OUT}-${VERSION} -tags netgo -ldflags='-extldflags \"-static\" -w -s -X main.version=${VERSION}' ${PKG}"

run: server
	./${OUT}

clean:
	-@del ${OUT} ${OUT}-* 2>nul

.PHONY: run server static vet lint

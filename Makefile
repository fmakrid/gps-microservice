OUT := gps-microservice
PKG := github.com/fmakrid/gps-microservice
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | findstr /V "\\vendor")
GO_FILES := $(shell dir /s /b *.go | findstr /V "\\vendor")


all: run

server:
	go build -v -o ${OUT} -ldflags="-X main.version=${VERSION}" ${PKG}

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

static: vet lint
	go build -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}" ${PKG}

run: server
	./${OUT}

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run server static vet lint

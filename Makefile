SHELL := /bin/bash
BASEDIR = $(shell pwd)

versionDir = "github.com/vixyninja/go-blocks/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Ho_Chi_Minh date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w \
-X ${versionDir}.Version=${gitTag} \
-X ${versionDir}.Commit=${gitCommit} \
-X ${versionDir}.BuiltAt=${buildDate}"

PROJECT_NAME := "github.com/vixyninja/go-blocks"
PKG := "$(PROJECT_NAME)"

export PATH        := $(shell go env GOPATH)/bin:$(PATH)
export GOPATH      := $(shell go env GOPATH)
export GO111MODULE := on

FOLDER ?= ./...

.PHONY: build test fmt vet tidy lint

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	@command -v golangci-lint >/dev/null 2>&1 || \
		(echo "golangci-lint not found, installing..."; \
		GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run $(FOLDER)

.PHONY: version
version:
	@echo "Package:      ${PKG}"
	@echo "Version:      $(gitTag)"
	@echo "Commit:       $(gitCommit)"
	@echo "BuiltAt:      $(buildDate)"
	@echo "TreeState:    $(gitTreeState)"

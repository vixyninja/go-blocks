FOLDER ?= ./...

.PHONY: build test fmt vet tidy lint

build:
	go build ./...

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	@command -v golangci-lint >/dev/null 2>&1 || \
		(echo "golangci-lint not found, installing..."; \
		GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run $(FOLDER)


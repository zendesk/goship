VERSION := $(shell git rev-parse HEAD)
all: build

prepare:
	mkdir build

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o build/goship_darwin_amd64
build_linux:
	GOOS=linux GOARCH=amd64 go build -o build/goship_linux_amd64

build: build_darwin build_linux

clean:
	rm build/*

test-lint:
	@echo
	@echo "==> Running linters <=="
	scripts/validate-go.sh

test-unit:
	@echo
	@echo "==> Running unit tests <=="
	go test -run . -race -cover ./...

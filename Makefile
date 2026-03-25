.PHONY: all build test clean install run

BINARY_NAME=gh-repo-visualize
MAIN_PACKAGE=./cmd/gh-repo-visualize

all: test build

build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

test:
	go test -v ./...

clean:
	rm -rf bin/

install: build
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/

run: build
	./bin/$(BINARY_NAME)

# Development
fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

# Cross-platform builds
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

build-all: build-linux build-darwin build-windows

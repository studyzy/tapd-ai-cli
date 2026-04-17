.PHONY: build test test-integration lint coverage clean fmt install

BINARY_NAME=tapd
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/studyzy/tapd-ai-cli/internal/cmd.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/tapd/

test:
	go test ./...
	cd sdk && go test ./...

test-integration:
	go test ./... -v -run "TestIntegration" -count=1

lint:
	go vet ./...
	goimports -l .

fmt:
	gofmt -w .
	goimports -w .

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	cd sdk && go test ./... -coverprofile=coverage.out
	cd sdk && go tool cover -func=coverage.out

install:
	go install $(LDFLAGS) ./cmd/tapd/

clean:
	rm -f $(BINARY_NAME) coverage.out

.PHONY: build test test-integration lint coverage clean fmt

BINARY_NAME=tapd

build:
	go build -o $(BINARY_NAME) ./cmd/tapd/

test:
	go test ./...

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

clean:
	rm -f $(BINARY_NAME) coverage.out

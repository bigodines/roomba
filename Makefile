export GO111MODULE=on

.PHONY: test
test:
	go test --cover --race

.PHONY: build
build:
	go build

.PHONY: lint
lint:
	golangci-lint run ./...

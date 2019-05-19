export GO111MODULE=on

.PHONY: test
test:
	$(GO11MODULE) go test --cover --race

.PHONY: build
build:
	$(GO11MODULE) go build

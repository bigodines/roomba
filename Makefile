export GO111MODULE=on

CMDS=$(filter-out internal, $(notdir $(wildcard cmd/*)))

.PHONY: test
test:
	go test --cover --race -v ./lib/*

.PHONY: build
build: $(CMDS)

$(CMDS):
	go build  -o ./bin/$@ cmd/$@/*.go

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run ./...

export GO111MODULE=on

CMDS=$(filter-out internal, $(notdir $(wildcard cmd/*)))

.PHONY: test
test:
	go test --cover --race

.PHONY: build
build: $(CMDS)

$(CMDS):
	go build  -o ./bin/$@ cmd/$@/*.go

.PHONY: lint
lint:
	golangci-lint run ./...

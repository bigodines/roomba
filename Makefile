export GO111MODULE=on

CMDS=$(filter-out internal, $(notdir $(wildcard cmd/*)))

### Local dev ----
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

### Docker -----

## TODO: build container for each GOOS
.PHONY: docker
docker:
	GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o ./bin/cron ./cmd/cron/*.go
	docker build -t olcolabs/roomba .

.PHONY: push
push:
	docker push olcolabs/roomba:latest

.PHONY: run
run:
	docker run --env GITHUB_TOKEN=$(GITHUB_TOKEN) --entrypoint=./cron olcolabs/roomba:latest

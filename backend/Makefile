.PHONY: all
all: generate test lint

.PHONY: generate
generate:
	cd ent && go generate

.PHONY: run
run:
	go run cmd/main.go

.PHONY: test
test:
	go test ./...

.PHONY: precommit
precommit: lint fmt imports

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: imports
imports:
	goimports -w .
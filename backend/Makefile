UNAME_S := $(shell uname -s)

ifeq ($(findstring MINGW,$(UNAME_S)),MINGW)
	PYTHON := python
else
	PYTHON := python3
endif

.PHONY: all
all: generate test lint

.PHONY: setup
setup:
ifeq ($(findstring MINGW,$(UNAME_S)),MINGW)
	@echo "Detected Windows (MinGW)"
	$(PYTHON) get-pip.py
	$(PYTHON) -m pip install --upgrade pip
	$(PYTHON) -m pip install pre-commit
else
	@echo "Detected Unix-like system (macOS/Linux)"
	brew install pre-commit
endif
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	npm install -g @commitlint/cli @commitlint/config-conventional
	pre-commit install
	pre-commit autoupdate --repo https://github.com/pre-commit/pre-commit-hooks

.PHONY: generate
generate:
	cd ent && go generate

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run --fix

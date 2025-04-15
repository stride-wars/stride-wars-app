setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

generate:
	cd internal/ent && go generate

run:
	go run cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run --fix


.PHONY: run dev seed build test test-integration test-cover lint fmt tidy docker-build

run:
	go run ./cmd/server

dev:
	go run github.com/air-verse/air@latest -c .air.toml

seed:
	go run ./cmd/server -seed

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	golangci-lint run

fmt:
	gofmt -l -w .
	goimports -w .

tidy:
	go mod tidy

docker-build:
	docker build -t penates-backend .

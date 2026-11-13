.PHONY: build run test test-integration test-coverage swagger lint up down clean

APP_NAME = recodations-api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -v ./cmd/... ./internal/... ./pkg/... ./test/...

test-short:
	go test -short -v ./cmd/... ./internal/... ./pkg/... ./test/...

test-coverage:
	go test -coverprofile=coverage.out ./cmd/... ./internal/... ./pkg/... ./test/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

swagger:
	swag init -g cmd/api/main.go -o docs --parseDependency

up:
	docker compose up -d

down:
	docker compose down

clean:
	rm -rf bin/ coverage.out coverage.html

.PHONY: build run test test-integration test-coverage swagger lint up down clean

APP_NAME = recodations-api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -v ./...

test-short:
	go test -short -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
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

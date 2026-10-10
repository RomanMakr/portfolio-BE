.PHONY: help tidy run test fmt db-up db-down db-reset db-logs

help:
	@echo "Targets: tidy run test fmt db-up db-down db-reset db-logs"

tidy:
	go mod tidy

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-reset:
	docker compose down -v
	docker compose up -d postgres

db-logs:
	docker compose logs -f postgres

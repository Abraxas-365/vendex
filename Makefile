.PHONY: dev up down migrate build test run

# Infrastructure
up:
	docker compose up -d

down:
	docker compose down

migrate:
	@for f in migrations/*.up.sql; do \
		echo "Running $$f..."; \
		psql "$(DATABASE_URL)" -f "$$f"; \
	done

# Development
dev:
	@echo "Starting backend..."
	go run ./cmd/...

run: build
	./bin/hada-commerce

build:
	go build -o bin/hada-commerce ./cmd/...

test:
	go test ./...

# Lint
vet:
	go vet ./...

check: vet test
	@echo "All checks passed."

.PHONY: dev up down backend frontend

# Infrastructure
up:
	docker compose up -d

down:
	docker compose down

# Backend
backend:
	cd backend && go run ./cmd/...

backend-build:
	cd backend && go build -o ../bin/hada-commerce ./cmd/...

backend-test:
	cd backend && go test ./...

# Frontend
frontend:
	cd frontend && bun run dev

frontend-build:
	cd frontend && bun run build

# Run everything
dev:
	@echo "Run 'make up' first for Postgres+Redis"
	@echo "Then in separate terminals: 'make backend' and 'make frontend'"

# Migrate
migrate:
	@for f in backend/migrations/*.up.sql; do \
		echo "Running $$f..."; \
		psql "$(DATABASE_URL)" -f "$$f"; \
	done

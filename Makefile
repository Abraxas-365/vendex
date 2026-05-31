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
	cd backend && go build -o ../bin/vendex ./cmd/...

backend-test:
	cd backend && go test ./...

# Frontend
frontend:
	cd frontend && bun run dev

frontend-build:
	cd frontend && bun run build

# Run everything in a tmux split
dev:
	@tmux kill-session -t vendex 2>/dev/null; true
	@tmux new-session -d -s vendex -n dev 'cd $(CURDIR)/backend && go run ./cmd/...'
	@tmux split-window -h -t vendex 'cd $(CURDIR)/frontend && bun run dev'
	@tmux select-pane -t vendex:dev.0
	@tmux switch-client -t vendex 2>/dev/null || tmux attach -t vendex

# Migrate
migrate:
	@for f in backend/migrations/*.up.sql; do \
		echo "Running $$f..."; \
		psql "$(DATABASE_URL)" -f "$$f"; \
	done

# Agent Presets
PRESET_REGISTRY ?= ghcr.io/abraxas-365/vendex-presets
PRESET_TAG ?= latest

preset-build-base:
	docker build -t vendex-preset-base:$(PRESET_TAG) deploy/presets/base/

preset-build-webdev: preset-build-base
	docker build -t $(PRESET_REGISTRY)/webdev:$(PRESET_TAG) deploy/presets/webdev/

preset-build-researcher: preset-build-base
	docker build -t $(PRESET_REGISTRY)/researcher:$(PRESET_TAG) deploy/presets/researcher/

preset-build: preset-build-webdev preset-build-researcher
	@echo "All preset images built."

preset-push:
	docker push $(PRESET_REGISTRY)/webdev:$(PRESET_TAG)
	docker push $(PRESET_REGISTRY)/researcher:$(PRESET_TAG)

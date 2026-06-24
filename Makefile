.PHONY: dev dev-deps test web-dev web-build migrate-up migrate-down bootstrap-admin release-check deploy-config deploy-deps deploy-migrate deploy-up deploy-down deploy-logs deploy-backup

DEPLOY_ENV ?= staging
ENV_FILE ?= deploy/.env.$(DEPLOY_ENV)
DEV_ENV_FILE ?= .env
GO_CACHE ?= $(CURDIR)/.cache/go-build
DEV_ENV = set -a; if [ -f $(DEV_ENV_FILE) ]; then . $(DEV_ENV_FILE); else . .env.example; fi; set +a;
COMPOSE = docker compose --env-file $(ENV_FILE) -f deploy/compose.yaml -f deploy/compose.$(DEPLOY_ENV).yaml

dev:
	$(DEV_ENV) cd src/server && go run ./cmd/server

dev-deps:
	docker compose -f deploy/docker-compose.yml up -d postgres redis

test:
	cd src/server && GOCACHE=$(GO_CACHE) go test ./...
	cd src/web && bun run api:check && bun run typecheck && bun run test

web-dev:
	cd src/web && bun run dev

web-build:
	cd src/web && bun run typecheck && bun run build

migrate-up:
	$(DEV_ENV) cd src/server && go run github.com/pressly/goose/v3/cmd/goose@v3.27.1 -dir migrations/$${DB_DRIVER:-postgres} $${DB_DRIVER:-postgres} "$$DB_DSN" up

migrate-down:
	$(DEV_ENV) cd src/server && go run github.com/pressly/goose/v3/cmd/goose@v3.27.1 -dir migrations/$${DB_DRIVER:-postgres} $${DB_DRIVER:-postgres} "$$DB_DSN" down

bootstrap-admin:
	$(DEV_ENV) cd src/server && go run ./cmd/bootstrap-admin

release-check: test web-build
	docker compose --env-file deploy/.env.staging.example -f deploy/compose.yaml -f deploy/compose.staging.yaml config --quiet
	docker compose --env-file deploy/.env.prod.example -f deploy/compose.yaml -f deploy/compose.prod.yaml config --quiet

deploy-config:
	$(COMPOSE) config --quiet

deploy-deps: deploy-config
	$(COMPOSE) up -d postgres redis

deploy-migrate:
	set -a; . $(ENV_FILE); set +a; cd src/server && go run github.com/pressly/goose/v3/cmd/goose@v3.27.1 -dir migrations/postgres postgres "postgres://secondadmin:$$POSTGRES_PASSWORD@127.0.0.1:$${POSTGRES_PORT:-5432}/secondadmin?sslmode=disable" up

deploy-up: deploy-config
	$(COMPOSE) up -d --build

deploy-down:
	$(COMPOSE) down

deploy-logs:
	$(COMPOSE) logs -f api web gateway

deploy-backup:
	mkdir -p backups
	$(COMPOSE) exec -T postgres pg_dump -U secondadmin -d secondadmin -Fc > backups/secondadmin-$$(date +%Y%m%d-%H%M%S).dump

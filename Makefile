.PHONY: certs up down up-backend down-backend up-dev down-dev dev seed build test test-backend test-integration test-frontend lint fmt clean

certs:
	@command -v mkcert >/dev/null 2>&1 || { echo "mkcert not found — install it first: https://github.com/FiloSottile/mkcert#installation"; exit 1; }
	mkcert -install
	mkdir -p certs
	mkcert -cert-file certs/localhost.pem -key-file certs/localhost-key.pem localhost 127.0.0.1 ::1

up:
	docker compose up --build

down:
	docker compose down

up-backend:
	docker compose -f docker-compose.backend.yml up --build

down-backend:
	docker compose -f docker-compose.backend.yml down

up-dev:
	docker compose -f docker-compose.dev.yml up -d

down-dev:
	docker compose -f docker-compose.dev.yml down

dev: up-dev
	@echo "Postgres is up on localhost:$${POSTGRES_PORT:-5432}."
	@test -f certs/localhost.pem || echo "No local TLS certs found — run 'make certs' first (see README's HTTPS section)."
	@echo "In separate terminals, run:"
	@echo "  cd backend && make dev   # Air hot-reload on :$${PORT:-8080}"
	@echo "  cd frontend && make dev  # Vite dev server with HMR"

seed:
	docker compose exec backend /server -seed

build:
	docker compose build

test: test-backend test-frontend

test-backend:
	$(MAKE) -C backend test

test-integration:
	$(MAKE) -C backend test-integration

test-frontend:
	$(MAKE) -C frontend test

lint:
	$(MAKE) -C backend lint
	$(MAKE) -C frontend lint

fmt:
	$(MAKE) -C backend fmt
	$(MAKE) -C frontend format

clean:
	docker compose down -v

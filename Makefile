.PHONY: up down up-backend down-backend seed build test test-backend test-integration test-frontend lint fmt clean

up:
	docker compose up --build

down:
	docker compose down

up-backend:
	docker compose -f docker-compose.backend.yml up --build

down-backend:
	docker compose -f docker-compose.backend.yml down

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

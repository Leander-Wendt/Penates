# Penates

An inventory management webapp with a Gin Gonic backend and a Vue.js frontend. Developed for the requirements of the student body of the Berliner Hochschule für Technik (BHT).

## Architecture

A monorepo with three services, wired together via Docker Compose:

```
Browser → frontend (Vue 3 SPA, served by nginx) → /api, /uploads reverse-proxied → backend (Gin REST API) → PostgreSQL
```

The frontend always calls relative `/api/...` and `/uploads/...` paths; nginx proxies those to the backend container, so no build-time API URL baking is needed. Authentication is JWT bearer-token based, following the [MDN HTTP authentication framework](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Authentication). `POST /api/v1/auth/login` is additionally rate-limited per client IP (`LOGIN_RATE_LIMIT_ATTEMPTS` per `LOGIN_RATE_LIMIT_WINDOW`, default 5/minute) to slow down brute-force credential guessing.

Three roles: **Admin**, **Logistics**, **Student**. Only Admin and Logistics can manage Items, Users, and Organisations; only Admin can delete a User. Several Item fields (manufacturer, serial number, inspection dates, purchase date, price, resolution number) are only ever returned to Admin/Logistics. A Student's API responses omit them entirely.

## Tech stack

- **Backend**: Go, [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io/), JWT auth (`golang-jwt`), bcrypt password hashing.
- **Frontend**: Vue 3, TypeScript, Vite, Vue Router, Pinia, Tailwind CSS, `vue-i18n` (English/German).
- **Database**: PostgreSQL 16.
- **Testing**: TDD throughout — Go `testing` + `testify` + `testcontainers-go` (backend), Vitest + `@vue/test-utils` + `vitest-axe` (frontend).

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose (v2+)
- Go 1.23+ (for local backend development outside Docker)
- Node 20+ (for local frontend development outside Docker)
- [mkcert](https://github.com/FiloSottile/mkcert#installation) (for locally-trusted HTTPS certs, see [HTTPS in local development](#https-in-local-development))

## Quickstart

```sh
cp .env.example .env
make certs   # generates locally-trusted TLS certs into certs/ (requires mkcert)
make up
```

This builds and starts Postgres, the backend, and the frontend. Once healthy, open `https://localhost:3000`. Log in with the seeded admin account (`SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` from `.env`, defaults `admin@penates.local` / `change-me-to-a-strong-password` — **change these before any real deployment**).

To also populate the database with dummy data for local development or a demo:

```sh
make seed
```

This runs against whichever stack is currently up (`make up` or `make up-backend`) by executing `/server -seed` inside the running `backend` container; it's idempotent, so running it again won't create duplicates. Outside Docker, use `cd backend && make seed` instead. See `backend/internal/database/seed.go` for the exact dummy accounts/items created (e.g. `student1@penates.local` / `student-pw-123`).

## Local development

**DB only, FE + BE natively with hot reload** (recommended for day-to-day development — Air rebuilds/restarts the backend on save, Vite gives the frontend HMR):

```sh
cp .env.example .env
make up-dev           # starts postgres via docker-compose.dev.yml, exposed on POSTGRES_PORT (default 5432)

# in a separate terminal
cd backend
cp .env.example .env
make dev              # Air hot-reload, rebuilds on every .go file save

# in a separate terminal
cd frontend
make install
make dev              # Vite dev server with HMR, proxies /api and /uploads to the backend
```

`make down-dev` stops Postgres again. `make dev` (from the repo root) is a shortcut that starts Postgres and prints the two commands above.

**Backend + DB only** (useful when developing just the frontend against a real, containerized backend — no backend hot reload):

```sh
cp .env.example .env
make up-backend      # starts postgres + backend via docker-compose.backend.yml
cd frontend
make install
make dev              # Vite dev server, proxies /api and /uploads to the backend
```

**Backend natively** (without Docker, against a local or containerized Postgres):

```sh
cd backend
cp .env.example .env
make run
```

**Frontend natively**:

```sh
cd frontend
make install
make dev
```

## HTTPS in local development

Everything serves over HTTPS using a locally-trusted certificate, so it matches how the app is actually deployed and browsers don't treat the login form or JWT cookies as running on an insecure origin.

```sh
make certs
```

This installs [mkcert](https://github.com/FiloSottile/mkcert#installation)'s local CA into your OS/browser trust stores (`mkcert -install`) and writes `certs/localhost.pem` / `certs/localhost-key.pem` (gitignored — regenerate them yourself, don't commit them). Run it once per machine; both the backend and frontend pick the files up automatically after that:

- **Docker** (`docker-compose.yml` and `docker-compose.backend.yml`) mounts `./certs` into both the `backend` and `frontend` containers. The backend serves HTTPS directly on its port; nginx terminates TLS for the frontend on `443`, so `FRONTEND_PORT` (default `3000`) now maps to HTTPS, e.g. `https://localhost:3000`.
- **Backend natively** reads `TLS_CERT_FILE` / `TLS_KEY_FILE` from `backend/.env` (defaulted to `../certs/localhost.pem` / `../certs/localhost-key.pem` in `backend/.env.example`) and serves HTTPS on `PORT` when both are set. Leave them unset to fall back to plain HTTP (e.g. for tests).
- **Frontend natively** (Vite): `vite.config.ts` automatically serves HTTPS if `certs/localhost.pem`/`localhost-key.pem` exist, falling back to HTTP otherwise so `vite build` never depends on them. The dev server proxies `/api` and `/uploads` to `VITE_API_PROXY_TARGET` (default `https://localhost:8080`).

If you skip `make certs`, everything still works over plain HTTP. Just update `CORS_ALLOWED_ORIGINS`, `VITE_API_PROXY_TARGET`, and the URLs below back to `http://`.

Because a native Vite dev server runs on port `5173` by default (not `3000`), make sure `CORS_ALLOWED_ORIGINS` on the backend includes `https://localhost:5173` when developing that way. See `.env.example` / `backend/.env.example`, which already list both `3000` and `5173`.

## Running tests

```sh
make test              # backend unit tests + frontend tests
make test-backend      # backend unit tests only (no Docker required)
make test-integration  # backend repository integration tests (testcontainers-go — requires a Docker daemon)
make test-frontend     # frontend unit/component tests (vitest + vitest-axe)
make lint               # backend (golangci-lint) + frontend (eslint) lint
```

## Accessing the database

Postgres isn't exposed to the host by `docker-compose.yml` or `docker-compose.backend.yml` (kept off the host network by default); `docker-compose.dev.yml` does expose it on `POSTGRES_PORT` so natively-run backend/tooling can reach it. To inspect the data, run `psql` inside the running `postgres` container:

```sh
docker compose exec postgres psql -U postgres -d penates
```

(credentials come from `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` in `.env`, `postgres` / `postgres` / `penates` by default). Once connected: `\dt` lists tables, `\d <table>` describes one, `\q` quits. Or run a one-off query without an interactive session:

```sh
docker compose exec postgres psql -U postgres -d penates -c "SELECT inventory_number, name, category FROM items;"
```

## Project structure

```
Penates/
├── docker-compose.yml               # full stack: postgres + backend + frontend
├── docker-compose.backend.yml       # postgres + backend only
├── docker-compose.dev.yml           # postgres only, port exposed for native FE + BE hot reload
├── Makefile                          # root orchestration, delegates into backend/ and frontend/
├── backend/                          # Gin API — cmd/, internal/{config,database,models,repository,service,handler,middleware,auth,dto,testutil}
└── frontend/                         # Vue 3 SPA — src/{views,components,stores,router,api,types,composables,locales}
```

## API overview

Base path `/api/v1` (plus `/healthz` and `/uploads/:filename`). All routes except `/healthz` and `/auth/login` require `Authorization: Bearer <token>`.

| Route                                  | Method                                      | Roles                                                  |
| -------------------------------------- | ------------------------------------------- | ------------------------------------------------------ |
| `/healthz`                             | GET                                         | public (also checks DB connectivity)                   |
| `/api/v1/auth/login`                   | POST                                        | public                                                 |
| `/api/v1/organisations`                | GET / POST                                  | any authenticated / Admin+Logistics                    |
| `/api/v1/organisations/:id`            | GET / PUT / DELETE                          | any / Admin+Logistics / Admin+Logistics                |
| `/api/v1/users/me`                     | GET                                         | any authenticated                                      |
| `/api/v1/users`                        | GET / POST                                  | Admin+Logistics                                        |
| `/api/v1/users/:id`                    | GET / PUT                                   | Admin+Logistics (a user may also GET their own record) |
| `/api/v1/users/:id`                    | DELETE                                      | Admin only                                             |
| `/api/v1/items`                        | GET (`?search=&category=&location=`) / POST | any authenticated / Admin+Logistics                    |
| `/api/v1/items/:inventoryNumber`       | GET / PUT / DELETE                          | any / Admin+Logistics / Admin+Logistics                |
| `/api/v1/items/:inventoryNumber/image` | POST / DELETE                               | Admin+Logistics                                        |
| `/api/v1/loan-requests`                | GET (role-filtered) / POST                  | any authenticated                                      |
| `/api/v1/loan-requests/:id`            | GET / PUT / DELETE                          | owner or Admin+Logistics                               |
| `/api/v1/loan-requests/:id/status`     | PATCH                                       | Admin+Logistics                                        |
| `/uploads/:filename`                   | GET                                         | any authenticated                                      |

## Data model

**Item** — primary key `inventoryNumber` (6-digit string): `name` (required), `description`, `imagePath` (set via the image upload endpoint), `category`, `location`, `amount` (default 1), `electricalAppliance` (boolean), `note`, and — visible to Admin/Logistics only — `manufacturer`, `serialNumber`, `lastTechnicalInspectionDate`, `lastElectricalInspectionDate`, `dateOfPurchase`, `price`, `resolutionNumber`.

**User** — `id`, `email` (required, unique), `password` (write-only, bcrypt-hashed), `role` (`admin` / `logistics` / `student`, required), `organisationId` (FK to an admin/logistics-managed `Organisation`).

**Organisation** — `id`, `name` (required, unique), `description`.

**LoanRequest** — `id`, `items` (the borrowed `Item`s), `requestingUserId`, `dateOfLending`, `dateOfReturn`, `locationOfItems`, `status` (`pending` / `approved` / `rejected` / `returned` / `cancelled`, default `pending`).

## Environment variables

| Variable                                                      | Used by                  | Default                                         | Purpose                                                                                                                                         |
| ------------------------------------------------------------- | ------------------------ | ----------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB`         | postgres, backend        | `postgres` / `postgres` / `penates`             | Database credentials                                                                                                                            |
| `BACKEND_PORT`                                                | root compose             | `8080`                                          | Backend host port (not exposed in the full-stack compose, only via `docker-compose.backend.yml`)                                                |
| `FRONTEND_PORT`                                               | root compose             | `3000`                                          | Host port the frontend is served on                                                                                                             |
| `POSTGRES_PORT`                                               | `docker-compose.dev.yml` | `5432`                                          | Host port Postgres is exposed on for native backend development                                                                                 |
| `JWT_SECRET`                                                  | backend                  | —                                               | Secret used to sign JWTs — set a real value before any real deployment                                                                          |
| `JWT_EXPIRY`                                                  | backend                  | `24h`                                           | Token lifetime                                                                                                                                  |
| `CORS_ALLOWED_ORIGINS`                                        | backend                  | `https://localhost:3000,https://localhost:5173` | Allowed CORS origins                                                                                                                            |
| `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` / `SEED_ADMIN_ORG` | backend                  | see `.env.example`                              | First Admin account, auto-created on startup if the `users` table is empty                                                                      |
| `MAX_UPLOAD_SIZE_MB`                                          | backend                  | `10`                                            | Max item image upload size                                                                                                                      |
| `LOGIN_RATE_LIMIT_ATTEMPTS`                                   | backend                  | `5`                                             | Max `/api/v1/auth/login` attempts allowed per client IP within `LOGIN_RATE_LIMIT_WINDOW`                                                        |
| `LOGIN_RATE_LIMIT_WINDOW`                                     | backend                  | `1m`                                            | Time window for the login rate limit, as a Go duration (e.g. `1m`, `30s`)                                                                       |
| `TLS_CERT_FILE` / `TLS_KEY_FILE`                              | backend                  | unset (Docker sets `/certs/...`)                | Paths to the TLS cert/key; serves HTTPS when both are set, plain HTTP otherwise — see [HTTPS in local development](#https-in-local-development) |

See `backend/.env.example` and `frontend/.env.example` for the per-service subsets used when running a service outside Docker.

## License

MIT — see [LICENSE](LICENSE).

# Penates

An inventory management webapp with a Gin Gonic backend and a Vue.js frontend. Developed for the requirements of the student body of the Berliner Hochschule für Technik (BHT).

## Architecture

A monorepo with three services, wired together via Docker Compose:

```
Browser → frontend (Vue 3 SPA, served by nginx) → /api, /uploads reverse-proxied → backend (Gin REST API) → PostgreSQL
```

The frontend always calls relative `/api/...` and `/uploads/...` paths; nginx proxies those to the backend container, so no build-time API URL baking is needed. Authentication is JWT bearer-token based, following the [MDN HTTP authentication framework](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Authentication): `401` (with a `WWW-Authenticate: Bearer` header) means the token is missing/invalid/expired, `403` means the token is valid but the user's role isn't permitted for that route.

Three roles: **Admin**, **Logistics**, **Student**. Only Admin and Logistics can manage Items, Users, and Organisations; only Admin can delete a User. Several Item fields (manufacturer, serial number, inspection dates, purchase date, price, resolution number) are only ever returned to Admin/Logistics — a Student's API responses omit them entirely.

## Tech stack

- **Backend**: Go, [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io/) (PostgreSQL driver), JWT auth (`golang-jwt`), bcrypt password hashing.
- **Frontend**: Vue 3, TypeScript, Vite, Vue Router, Pinia, Tailwind CSS, `vue-i18n` (English/German).
- **Database**: PostgreSQL 16.
- **Testing**: TDD throughout — Go `testing` + `testify` + `testcontainers-go` (backend), Vitest + `@vue/test-utils` + `vitest-axe` (frontend).

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose (v2+)
- Go 1.23+ (for local backend development outside Docker)
- Node 20+ (for local frontend development outside Docker)
- `make`

## Quickstart

```sh
cp .env.example .env
make up
```

This builds and starts Postgres, the backend, and the frontend. Once healthy, open `http://localhost:3000` (or whatever `FRONTEND_PORT` you set in `.env`). Log in with the seeded admin account (`SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` from `.env`, defaults `admin@penates.local` / `change-me-to-a-strong-password` — **change these before any real deployment**).

To also populate the database with dummy data for local development or a demo:

```sh
make seed
```

This runs against whichever stack is currently up (`make up` or `make up-backend`) by executing `/server -seed` inside the running `backend` container; it's idempotent, so running it again won't create duplicates. Outside Docker, use `cd backend && make seed` instead. See `backend/internal/database/seed.go` for the exact dummy accounts/items created (e.g. `student1@penates.local` / `student-pw-123`).

## Local development

**Backend + DB only** (useful when developing the frontend with `npm run dev` against a real backend):

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

## Running tests

```sh
make test              # backend unit tests + frontend tests
make test-backend      # backend unit tests only (no Docker required)
make test-integration  # backend repository integration tests (testcontainers-go — requires a Docker daemon)
make test-frontend     # frontend unit/component tests (vitest + vitest-axe)
make lint               # backend (golangci-lint) + frontend (eslint) lint
```

## Project structure

```
Penates/
├── docker-compose.yml               # full stack: postgres + backend + frontend
├── docker-compose.backend.yml       # postgres + backend only
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

| Variable                                                      | Used by           | Default                             | Purpose                                                                                          |
| ------------------------------------------------------------- | ----------------- | ----------------------------------- | ------------------------------------------------------------------------------------------------ |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB`         | postgres, backend | `postgres` / `postgres` / `penates` | Database credentials                                                                             |
| `BACKEND_PORT`                                                | root compose      | `8080`                              | Backend host port (not exposed in the full-stack compose, only via `docker-compose.backend.yml`) |
| `FRONTEND_PORT`                                               | root compose      | `3000`                              | Host port the frontend is served on                                                              |
| `JWT_SECRET`                                                  | backend           | —                                   | Secret used to sign JWTs — set a real value before any real deployment                           |
| `JWT_EXPIRY`                                                  | backend           | `24h`                               | Token lifetime                                                                                   |
| `CORS_ALLOWED_ORIGINS`                                        | backend           | `http://localhost:3000`             | Allowed CORS origins                                                                             |
| `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` / `SEED_ADMIN_ORG` | backend           | see `.env.example`                  | First Admin account, auto-created on startup if the `users` table is empty                       |
| `MAX_UPLOAD_SIZE_MB`                                          | backend           | `10`                                | Max item image upload size                                                                       |

See `backend/.env.example` and `frontend/.env.example` for the per-service subsets used when running a service outside Docker.

## License

MIT — see [LICENSE](LICENSE).

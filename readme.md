# Transactions API

Go service that manages bank accounts and transactions. It exposes REST endpoints for creating accounts, querying accounts, and recording transactions with simple validation. The server runs on Echo, persists data in PostgreSQL, and applies migrations automatically on startup.

## Requirements

- Go 1.25+
- PostgreSQL 14+ (Docker Compose spins up Postgres 17)
- Docker and Docker Compose (for containerized runs)
- `make` (optional, for the coverage helper)

## Configuration

Runtime settings come from environment variables. A starter `.env` is included:

```
ENVIRONMENT=local
DB_HOST=127.0.0.1
DB_USER=postgres
DB_PASS=postgres
DB_PORT=5432
DB_SSL_MODE=disable
DB_NAME=transactions
DB_POOL_MINIMUM=2
DB_POOL_MAXIMUM=10
PORT=5500
```

Notes:
- The application auto-runs migrations against `DB_NAME` on startup.
- Tests create and migrate a `<DB_NAME>_test` database; the configured user must be allowed to create databases.

## Run the project

### Option 1: Docker Compose (recommended)

```
docker compose up --build
```

Services:
- API at `http://localhost:5500`
- PostgreSQL at `localhost:5432` (user/password: `postgres`/`postgres`, db: `transactions`)
- Swagger UI at `http://localhost:8080` (serving `openapi/openapi.json`)

### Option 2: Local Go run

1. Start a PostgreSQL instance matching `.env` (database `transactions` must exist).
2. Install dependencies: `go mod download`
3. Launch the API:
   ```
   go run ./cmd
   ```

The server listens on `:$PORT` (default `5500`) and applies migrations from `./migrations`.

## Running tests

Prerequisites: a PostgreSQL instance reachable with the `.env` settings, and a user allowed to create databases.

```
go test ./...
```

The test harness will:
- Load variables from `.env`
- Create and migrate `<DB_NAME>_test`
- Boot the API on the configured `PORT` to exercise HTTP endpoints

Coverage helper:
```
make cover
```
This generates `coverage.out` and opens the HTML report.

## Distributed tracing (Jaeger + OpenTelemetry)

- The service is instrumented with OpenTelemetry (Echo middleware + GORM plugin) and exports spans via OTLP/HTTP.
- In Docker Compose, Jaeger all-in-one runs as `jaeger`, and the API sends traces to `http://jaeger:4318`. The UI is available at `http://localhost:16686`.
- Key OTEL env vars (set in `docker-compose.yml`):
  - `OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4318`
- `OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf`
- `OTEL_TRACES_EXPORTER=otlp`
- Local (without Docker): point `OTEL_EXPORTER_OTLP_ENDPOINT` to your collector/Jaeger OTLP endpoint, then start the app. Keep `propagation.TraceContext` enabled by default to preserve cross-service context.
- Every HTTP request generates a span with request/response metadata; database calls emit spans via the GORM OpenTelemetry plugin. Use Jaeger’s “Search” to filter by service `transactions.api` and inspect end-to-end flows.

## Architecture (high level)

- **Transport**: Echo HTTP server with middleware for request logging/observability. Routes live in `internal/api`, grouped under `/api/v1`.
- **Validation & binding**: `internal/lib/rest` helpers bind/validate requests; custom validator in `internal/lib/validator`.
- **Domain/services**: Business logic in `internal/services` works with domain models from `internal/domain`.
- **Persistence**: Repositories in `internal/repositories` use GORM over PostgreSQL (`internal/lib/postgres`), with migrations in `/migrations` applied at startup.
- **Observability**: OpenTelemetry spans emitted for HTTP (middleware) and DB (GORM OTEL plugin); logs via `internal/lib/logging`.
- **Composition**: `cmd/main.go` wires env loading, tracing exporter, DB connection/pool, migrations, services, and routes.

## API surface

- Accounts: `POST /api/v1/accounts`, `GET /api/v1/accounts/:account_id`
- Transactions: `POST /api/v1/transactions`

OpenAPI spec lives in `openapi/openapi.json`; Swagger UI is available via Docker Compose (`http://localhost:8080`).

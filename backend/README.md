# Backend

Go REST API described in `../architecture/api-design.md`. Currently
implements Groups and Activity Types (§3.1–3.2); Planning, Recording, and
Summary endpoints are not built yet.

## Requirements

- Go 1.25+ (the module pins its own toolchain via `go.mod`/`go.sum`, so
  `go build`/`go test` fetch it automatically if needed).
- A Postgres database with `../migrations/0001_init.sql` applied.

## Environment variables

| Variable | Required | Notes |
|---|---|---|
| `DATABASE_URL` | yes | Postgres connection string (Neon's **direct**, non-pooled string — see `architecture/api-design.md` §6). |
| `CLERK_JWKS_URL` | yes | Clerk's JWKS endpoint, e.g. `https://<your-clerk-domain>/.well-known/jwks.json`. |
| `CLERK_WEBHOOK_SECRET` | yes | Clerk webhook signing secret (`whsec_...`), for `POST /webhooks/clerk`. |
| `PORT` | no | Defaults to `8080`. |

Clerk's default session token doesn't include an email claim — add a
custom claim named `email` to the session token template in the Clerk
dashboard if you want `internal/auth`'s JWT-path user upsert to have an
email on first login (the webhook path always has it).

## Running locally

```sh
# apply the schema to a local/dev Postgres
psql "$DATABASE_URL" -f ../migrations/0001_init.sql

go run ./cmd/api
```

## Testing

```sh
go build ./...
go vet ./...
gofmt -l .            # should print nothing

# unit tests (no database needed)
go test ./internal/auth/...

# full suite, including HTTP+Postgres integration tests
# (integration tests in internal/httpapi skip themselves if DATABASE_URL is unset)
DATABASE_URL="postgres://postgres:postgres@localhost:5432/tm_dev?sslmode=disable" go test ./...
```

## Regenerating sqlc code

Query files live in `db/query/*.sql`; generated code in
`internal/db/sqlcgen/` is committed so `go build` works without the sqlc
binary. After changing a query or `../migrations`, regenerate with:

```sh
go tool sqlc generate
```

(`sqlc` is a pinned tool dependency via `go.mod`'s `tool` directive — Go
1.24+'s `go tool` runs it without a separate global install.)

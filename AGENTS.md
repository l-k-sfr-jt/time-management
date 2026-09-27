# AGENTS.md

Instructions for AI coding agents (and a quick-start for humans) working in
this repository.

## Project

A personal time management app: track activities during the day, group
them by topic (Work/Sport/Relax/Study/Hobbies/...), plan ahead (including
recurring plans), and review planned-vs-actual time in weekly/monthly
summaries.

## Current status

The Go backend (`backend/`) implements Groups and Activity Types
(`architecture/api-design.md` §3.1–3.2): Clerk JWT auth, the Clerk
webhook, and Postgres access via pgx/sqlc. Planning, Recording, and
Summary endpoints, and the TanStack Start frontend, are not built yet —
check `architecture/api-design.md` for what's specified vs. implemented
before assuming an endpoint exists.

## Where things live

- `docs/requirements.md` — product functional/non-functional requirements.
  Read this before changing behavior — it's the source of truth for *what*
  the app does.
- `architecture/` — technical design (source of truth for *how* it's
  built):
  - `data-model.md` — entities and Postgres schema.
  - `data-flows.md` — sequence diagrams for the key flows.
  - `api-design.md` — Go backend API (endpoints, high-level design, deep
    dives on the harder trade-offs).
- `migrations/` — sequentially numbered SQL migrations
  (`0001_init.sql`, ...). Keep them in sync with `architecture/data-model.md`
  whenever the schema changes — update both in the same commit.

## Stack (see `docs/requirements.md` §9 for rationale)

- Backend: Go, REST API.
- Database: Postgres (Neon).
- Frontend: TanStack Start (React).
- Auth: Clerk.

## Working in this repo

- Before implementing a feature, check `architecture/api-design.md` for
  the endpoint contract and `architecture/data-model.md` for the schema —
  build to match them. If the right implementation diverges from what's
  documented, update the doc in the same change rather than letting them
  drift apart.
- Schema changes: add a new numbered file under `migrations/` (never edit
  an already-applied migration), and update `architecture/data-model.md`
  to match.
- Commit message format: see `CONTRIBUTING.md`.

## Backend build/lint/test

Run from `backend/` (see `backend/README.md` for env vars and full detail):

```sh
go build ./...
go vet ./...
gofmt -l .                     # should print nothing
go test ./internal/auth/...    # unit tests, no database needed
DATABASE_URL="postgres://postgres:postgres@localhost:5432/tm_dev?sslmode=disable" go test ./...  # full suite
```

Query files live in `backend/db/query/*.sql`; after changing one or
`migrations/`, regenerate the committed sqlc code with `go tool sqlc
generate` (run from `backend/`).

No frontend project exists yet — add its commands here once TanStack
Start is scaffolded.

# AGENTS.md

Instructions for AI coding agents (and a quick-start for humans) working in
this repository.

## Project

A personal time management app: track activities during the day, group
them by topic (Work/Sport/Relax/Study/Hobbies/...), plan ahead (including
recurring plans), and review planned-vs-actual time in weekly/monthly
summaries.

## Current status

**Design phase — no application code yet.** The repo currently holds
requirements and architecture docs plus one SQL migration. Do not assume
a running app, package manifest, or test suite exists; check before
referencing one.

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
- No build/lint/test commands exist yet. Once a backend or frontend
  project is scaffolded, add the actual commands here (e.g. `go build
  ./...`, `go test ./...`, the frontend's `npm run` scripts) so future
  agents don't have to rediscover them.

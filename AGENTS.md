# AGENTS.md

Instructions for AI coding agents (and a quick-start for humans) working in
this repository.

## Project

A personal time management app: track activities during the day, group
them by topic (Work/Sport/Relax/Study/Hobbies/...), plan ahead (including
recurring plans), and review planned-vs-actual time in weekly/monthly
summaries.

## Current status

**This branch (`learn/go-backend-fundamentals`) intentionally has no
`backend/` code.** It's a guided-learning rebuild: the repo owner is new
to Go and is rebuilding the backend themselves, one concept at a time,
following `lessons/`. A previously working implementation (Groups +
Activity Types, Clerk auth, pgx/sqlc, tests — see
`architecture/api-design.md` §3.1–3.2) exists for reference on branch
`claude/time-management-app-requirements-dgcxay`, but should not be
copied wholesale into this branch — the point is writing it from scratch.

## Lessons

- `lessons/README.md` — the curriculum: what each lesson covers and in
  what order. Start there.
- An AI agent asked to "help with the next lesson" should teach and
  review, not write the solution unprompted — see
  `lessons/README.md`'s "How sessions work" for the expected mode
  (explain the concept, assign a task, let the human write the code,
  review it, only show a reference solution if asked or the human is
  stuck after a real attempt).

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

No `backend/` module exists yet on this branch — lesson 01 creates it.
Once it exists, add the real commands here (they'll look like the ones
on `claude/time-management-app-requirements-dgcxay`: `go build ./...`,
`go vet ./...`, `gofmt -l .`, `go test ./...`) so this stays accurate.

No frontend project exists yet either.

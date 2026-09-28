# 11 — Postgres from Go

## Goal

Connect the Go service to a real Postgres database.

## Concepts

- Why `pgx` over `lib/pq` (maintenance-only) — recap of our earlier chat
  about the driver choice.
- `pgxpool.Pool` — a connection pool, not a single connection; why a
  long-running service wants pooling.
- Neon's direct vs. pooled (PgBouncer) connection strings, and why this
  service uses the direct one (see `architecture/api-design.md` §6).
- Applying `migrations/0001_init.sql` to a local Postgres with `psql`, so
  there's something to actually query against.

*(Task and checkpoint written when we start this lesson.)*

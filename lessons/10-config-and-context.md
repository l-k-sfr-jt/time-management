# 10 — Config & context

## Goal

Build `internal/config` (reading settings from environment variables)
and get an introduction to `context.Context`, which flows through nearly
every function signature from here on.

## Concepts

- `os.Getenv`, validating required config at startup rather than
  failing deep inside a request handler later.
- `context.Context` — carries request-scoped values (deadlines,
  cancellation signals, and in our case the authenticated user id) through
  a call chain. Convention: it's always the first parameter, named `ctx`.
- Why nearly every pgx/net-http function takes a `ctx context.Context` —
  it's how a cancelled request (client disconnected) or a timeout
  propagates down to "stop this database query."

*(Task and checkpoint written when we start this lesson.)*

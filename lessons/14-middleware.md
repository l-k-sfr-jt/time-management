# 14 — Middleware

## Goal

Learn Go's middleware pattern and use it for something real: a
panic-recovery middleware (from our earlier chat about chi — the one
piece worth having regardless of router choice).

## Concepts

- The pattern: `func(http.Handler) http.Handler` — a function that wraps
  one handler and returns another. This is Go's equivalent of Express's
  `(req, res, next) => {...}` middleware, just expressed through function
  composition instead of a `next()` callback.
- `defer` and `recover()` — how Go catches a panic (its closest thing to
  an exception, reserved for programmer errors, not normal control flow)
  so one bad request doesn't crash the whole server.
- Composing multiple middlewares around a handler/mux.

*(Task and checkpoint written when we start this lesson.)*

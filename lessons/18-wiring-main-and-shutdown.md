# 18 — Wiring main.go & graceful shutdown

## Goal

Assemble everything into `cmd/api/main.go`, and handle shutdown properly
instead of just killing the process.

## Concepts

- Goroutines (`go func() {...}()`) — Go's lightweight concurrency
  primitive; running the HTTP server on one while the main goroutine
  waits for a shutdown signal.
- Channels — how goroutines communicate; a small `chan error` here to
  learn about the server's own failure back to `main`.
- `signal.NotifyContext` — turning an OS signal (Ctrl-C, `SIGTERM`) into
  a `context.Context` that gets cancelled, tying back into lesson 10.
- `server.Shutdown(ctx)` — Go's built-in graceful shutdown: stop
  accepting new connections, let in-flight ones finish, with a timeout.

*(Task and checkpoint written when we start this lesson.)*

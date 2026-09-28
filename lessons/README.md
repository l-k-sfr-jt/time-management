# Learning Go by rebuilding this backend

This is a guided curriculum: we rebuild the Go backend from
`architecture/api-design.md`, one concept at a time, and **you** write
the code. It replaces the previously-generated implementation (still on
branch `claude/time-management-app-requirements-dgcxay` if you ever want
to compare notes), which isn't reused here on purpose.

## How sessions work

For each lesson:

1. I explain the Go concept(s) it introduces and what we're building.
2. I give you a concrete task — usually "make this file do X."
3. **You write the code.** Paste your attempt back as a code block in
   chat.
4. I review it: point out bugs, unidiomatic patterns, or missing pieces —
   without rewriting it for you. If something's off, you get another
   attempt. I'll only write the fix myself if you ask me to, or after a
   couple of genuine attempts where you're stuck.
5. Once it's right, I save your code to the actual file (as you wrote
   it, or with whatever specific fix you asked for) and we run the
   checkpoint (`go build`, `go test`, a `curl`, etc.) together to prove
   it works.
6. Commit, then move to the next lesson.

We can go faster through anything that feels familiar and slower through
anything that doesn't — this order is a default, not a rule.

## Prerequisites

None beyond general programming experience — that's the point. Each
lesson introduces the Go-specific parts it needs.

## Syllabus

### Part 1 — Go language fundamentals

| # | Lesson | Concepts |
|---|---|---|
| 01 | [Hello, Go](./01-hello-go.md) | modules (`go.mod`), `package main`, `func main`, `go run`/`build`/`fmt`/`vet` |
| 02 | [Types, control flow, functions](./02-types-control-flow-functions.md) | variables, basic types, `if`/`for`/`switch` (no `while`), functions, multiple return values |
| 03 | [Structs, methods, pointers](./03-structs-methods-pointers.md) | `struct`, methods with value vs. pointer receivers, when Go copies vs. shares |
| 04 | [Interfaces & errors](./04-interfaces-and-errors.md) | implicit interfaces, the `error` type, `errors.Is`/`errors.As`, wrapping with `%w` |
| 05 | [Slices & maps](./05-slices-maps.md) | slices vs. arrays, `append`, maps, zero values, `nil` |
| 06 | [Packages & project layout](./06-packages-and-project-layout.md) | exported vs. unexported names, `internal/`, how Go resolves imports |
| 07 | [Testing in Go](./07-testing-in-go.md) | `go test`, table-driven tests, `t.Run`, `t.Fatalf` vs `t.Errorf` |

### Part 2 — Building the HTTP service

| # | Lesson | Concepts | Builds |
|---|---|---|---|
| 08 | [net/http fundamentals](./08-http-fundamentals.md) | `http.Handler`, `ResponseWriter`, `*Request`, starting a server | `GET /healthz` |
| 09 | [Routing & JSON](./09-routing-and-json.md) | Go 1.22+ `ServeMux` method patterns, path params, `encoding/json` | request/response types |
| 10 | [Config & context](./10-config-and-context.md) | reading env vars, `context.Context` basics | `internal/config` |
| 11 | [Postgres from Go](./11-postgres-and-pgx.md) | `pgx`/`pgxpool`, connection pooling, why not `lib/pq` | `internal/db` |
| 12 | [sqlc](./12-sqlc.md) | generating typed Go from SQL, why not an ORM | `db/query/*.sql` |
| 13 | [Groups endpoint](./13-groups-endpoint.md) | putting 08–12 together end to end | full Groups CRUD |
| 14 | [Middleware](./14-middleware.md) | `func(http.Handler) http.Handler`, composing middleware | logging + panic recovery |
| 15 | [Activity Types endpoint](./15-activity-types-endpoint.md) | applying 08–14 with less hand-holding | full Activity Types CRUD |
| 16 | [JWT auth (Clerk)](./16-jwt-auth.md) | JWKS, RSA signature verification, request context values | auth middleware |
| 17 | [Webhooks & HMAC](./17-webhooks-hmac.md) | `crypto/hmac`, constant-time comparison, why it matters | Clerk webhook handler |
| 18 | [Wiring main.go](./18-wiring-main-and-shutdown.md) | goroutines, channels, `signal.NotifyContext`, graceful shutdown | `cmd/api/main.go` |
| 19 | [Integration testing](./19-integration-testing.md) | `httptest`, testing against a real Postgres, test skips | full test suite |

Each lesson file (once we reach it) has: **Goal**, **Concepts**, **Task**,
and **Checkpoint** (how we prove it works). Lessons 02–19 currently only
have Goal/Concepts filled in as a roadmap — I'll flesh out the Task and
Checkpoint right before we start each one, so it reflects where the code
actually stands by then.

Start with [lesson 01](./01-hello-go.md).

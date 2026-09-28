# 12 — sqlc

## Goal

Write your first real SQL query and generate type-safe Go code from it,
instead of hand-writing scan/bind code (or reaching for an ORM).

## Concepts

- Why not an ORM: recurrence-rule and time-range queries stay explicit
  and reviewable as plain SQL, matching the decision already recorded in
  `architecture/api-design.md`.
- `sqlc`: point it at your schema (`migrations/`) and a `.sql` file with
  named queries (`-- name: CreateGroup :one`); it generates a typed Go
  function + params struct per query.
- Go 1.24+'s `tool` directive in `go.mod` — pinning `sqlc` as a versioned
  tool dependency (`go get -tool`, `go tool sqlc generate`) instead of a
  separately-installed global binary.
- Reading generated code — it's meant to be read, not treated as a black
  box.

*(Task and checkpoint written when we start this lesson.)*

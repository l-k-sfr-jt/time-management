# 19 — Integration testing

## Goal

Write tests that exercise the real stack — HTTP layer through to
Postgres — building on the `go test` basics from lesson 07.

## Concepts

- `net/http/httptest` — spinning up a real (in-process) HTTP server for
  tests, and a fake JWKS server to sign test tokens against, so auth is
  tested for real rather than bypassed.
- Test skip pattern: `t.Skip(...)` when `DATABASE_URL` isn't set, so
  `go test ./...` still passes in an environment without Postgres (like
  CI without a database configured) instead of failing loudly.
- Testing a full lifecycle (create → list → conflict → archive) in one
  test vs. many small isolated tests — trade-offs of each.

*(Task and checkpoint written when we start this lesson. This is the
last lesson in the current curriculum — Planning/Recording/Summary
endpoints would extend it further once we get here.)*

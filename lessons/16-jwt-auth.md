# 16 — JWT auth (Clerk)

## Goal

Build the Clerk JWT verification middleware — the most involved piece of
the backend so far.

## Concepts

- JWTs: header (which key/algorithm signed it) + claims (payload) +
  signature; what "verifying" actually checks.
- JWKS — why the signing key isn't hardcoded: Clerk publishes its public
  keys at a well-known URL, and rotates them.
- A small in-memory cache with on-demand refresh (fetch once, reuse,
  refetch if a `kid` isn't found) — a common pattern anywhere you'd
  otherwise hit a network call per request.
- `context.Context` (from lesson 10) put to real use: attaching the
  resolved user id to the request context so handlers downstream can read
  it without it being passed explicitly through every function signature.
- `sync.RWMutex` — protecting the key cache from concurrent access (this
  is Go's answer to a race condition: multiple requests hitting the cache
  at once).

*(Task and checkpoint written when we start this lesson.)*

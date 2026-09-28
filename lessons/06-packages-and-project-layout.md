# 06 — Packages & project layout

## Goal

Understand how Go organizes code across files/folders, so the
`internal/config`, `internal/db`, `internal/httpapi` structure we're
about to build makes sense rather than feeling arbitrary.

## Concepts

- One package per directory; all files in a directory share the same
  `package` name and can see each other's unexported names directly
  (no imports needed within a package).
- **Exported vs. unexported**: capitalized name (`Group`) is visible
  outside the package; lowercase (`group`) is not. This *is* Go's
  access-control mechanism — there's no `private`/`public` keyword.
- The `internal/` convention: any package under a path containing
  `internal/` can only be imported by code inside that same module tree
  — Go's compiler enforces this, not just a naming convention.
- Import paths are based on the module path from `go.mod` +
  the folder path — how `github.com/l-k-sfr-jt/time-management/backend/internal/config`
  resolves to `backend/internal/config/`.

*(Task and checkpoint written when we start this lesson.)*

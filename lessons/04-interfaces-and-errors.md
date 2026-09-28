# 04 — Interfaces & errors

## Goal

Learn Go's two most distinctive idioms: implicit interfaces, and
error-values-instead-of-exceptions.

## Concepts

- **Implicit interfaces** — a type satisfies an interface just by having
  the right methods; there's no `implements` keyword. This is how
  `http.Handler` works (anything with a `ServeHTTP(w, r)` method *is* a
  handler).
- **`error` is just a value.** Go has no exceptions for normal error
  handling — functions return `(result, error)` and callers check `if err
  != nil`. Verbose compared to try/catch, but nothing is ever silently
  swallowed.
- Creating errors: `errors.New`, `fmt.Errorf`.
- **Wrapping**: `fmt.Errorf("doing X: %w", err)` preserves the original
  error so callers can unwrap it — `errors.Is`/`errors.As` — instead of
  losing context the way a re-thrown generic exception often does.
- Where this shows up in the real code: every `sqlc`-generated query
  returns `(Row, error)`, and our handlers check `errors.Is(err,
  pgx.ErrNoRows)` to turn "not found" into a 404.

*(Task and checkpoint written when we start this lesson.)*

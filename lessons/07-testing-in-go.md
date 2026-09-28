# 07 — Testing in Go

## Goal

Learn Go's built-in testing tools before we rely on them constantly from
lesson 13 onward.

## Concepts

- Testing needs no external framework: `*_test.go` files, functions
  named `TestXxx(t *testing.T)`, run with `go test ./...`.
- `t.Errorf` (report failure, keep running) vs. `t.Fatalf` (report and
  stop this test immediately).
- **Table-driven tests** — Go's idiomatic way to test many
  input/output pairs without repeating the test body:
  ```go
  cases := []struct{ in, want int }{{1, 2}, {2, 3}}
  for _, c := range cases {
      if got := increment(c.in); got != c.want {
          t.Errorf("increment(%d) = %d, want %d", c.in, got, c.want)
      }
  }
  ```
- `t.Run("name", func(t *testing.T) {...})` for subtests (shows up
  individually in `go test -v` output).
- `t.Skip(...)` — how the real integration tests (lesson 19) skip
  themselves when there's no database configured.

*(Task and checkpoint written when we start this lesson.)*

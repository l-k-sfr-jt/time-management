# 01 — Hello, Go

## Goal

Get a Go module running in `backend/` and understand the four commands
you'll use constantly: `go run`, `go build`, `go fmt`, `go vet`.

## Concepts

- **Modules.** A Go project is a *module*, declared by a `go.mod` file at
  its root. It names the module (an import path — for a private project
  like this, we just use the GitHub path it lives at) and the minimum Go
  version it needs. There's no `package.json`-style dependency list to
  hand-edit; `go.mod` grows automatically as you `go get` things later.
- **`package main` and `func main`.** Every Go file starts with a
  `package` declaration. A package named `main`, containing a function
  named `main`, is what makes a program *executable* (as opposed to a
  *library* package other code imports). `func main()` is the entry
  point — like `main()` in C, or the top-level code in a Node script.
- **No unused imports or variables, ever.** Go's compiler refuses to
  build code with an imported package you don't use, or a local variable
  you declare and never read. This isn't a linter suggestion — it's a
  compile error. Mildly annoying at first, genuinely useful once you're
  used to it (dead imports/variables can't silently accumulate).
- **The toolchain:**
  - `go run <path>` — compile and run in one step, for quick iteration.
  - `go build <path>` — compile to a binary, without running it.
  - `go fmt` (or `gofmt`) — the *canonical* formatter. There is
    deliberately no debate about brace placement or tabs vs. spaces in
    Go — `gofmt` decides, everyone runs it, done.
  - `go vet` — a static checker for common mistakes (e.g. a `Printf`
    call whose format string doesn't match its arguments).

## Task

1. Inside `backend/`, initialize a module:
   ```sh
   cd backend
   go mod init github.com/l-k-sfr-jt/time-management/backend
   ```
   This creates `backend/go.mod`. Open it and read it — it's short.

2. Create `backend/cmd/api/main.go` (the `cmd/api` path matches where the
   real service will live — see `architecture/api-design.md`). Write a
   program that:
   - Declares `package main`.
   - Imports the standard library `fmt` package.
   - Has a `func main()` that prints `Time management API booting...` to
     the console.

   Don't copy this from memory of other languages — think through: what
   does the file need at the top (package + import), and what's the
   single line inside `main` that prints a string?

3. Paste your `main.go` here and we'll review it together before running
   anything.

## Checkpoint

Once we're happy with the code, run it:

```sh
go run ./cmd/api
```

You should see your message printed. Then try:

```sh
go build ./cmd/api    # compiles to a binary, doesn't run it — where did the binary go?
gofmt -l .             # prints filenames that aren't formatted correctly (should be empty)
go vet ./...           # should report nothing
```

We'll commit `backend/go.mod` and `backend/cmd/api/main.go` once this all
passes, then move to [lesson 02](./02-types-control-flow-functions.md).

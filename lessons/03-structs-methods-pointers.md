# 03 — Structs, methods, pointers

## Goal

Understand Go's replacement for classes: structs plus methods, and the
value-vs-pointer distinction that trips up almost everyone coming from a
garbage-collected OO language.

## Concepts

- `struct` — a plain data type (no inheritance). This is what `Group`,
  `ActivityType`, etc. will be later.
- Methods: `func (g Group) String() string { ... }` — a function attached
  to a type via a *receiver*.
- **Value receiver vs. pointer receiver** — `func (g Group) Foo()` gets a
  *copy* of `g`; `func (g *Group) Foo()` gets a pointer to the original.
  This matters a lot once structs get large or need to be mutated.
- Pointers (`*T`, `&x`, `*p`) without the danger: Go has pointers but no
  pointer arithmetic, and the garbage collector means you never manually
  free anything.
- Zero values: a struct's fields start at their type's zero value (`0`,
  `""`, `false`, `nil`) — there's no "uninitialized" the way some
  languages have.

*(Task and checkpoint written when we start this lesson.)*

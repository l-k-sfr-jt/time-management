# 02 — Types, control flow, functions

## Goal

Get comfortable writing basic Go logic: variables, the three control-flow
constructs, and functions with multiple return values (which you'll see
*everywhere* in Go — it's how error handling works).

## Concepts

- Variable declaration: `var x int`, `x := 5` (short declaration, type
  inferred), and why Go is statically typed but rarely makes you write
  the type out.
- Basic types: `int`, `float64`, `string`, `bool`, and Go's explicitness
  about numeric types (no silent int→float conversions).
- `if`/`for`/`switch` — Go only has `for` (no `while`/`do-while`; a bare
  `for {}` is your infinite loop and `for condition {}` is your `while`).
- Functions: parameters, multiple return values (`func divide(a, b int)
  (int, error)`), named returns, and why this shape is everywhere once we
  get to error handling in lesson 04.

*(Task and checkpoint written when we start this lesson.)*

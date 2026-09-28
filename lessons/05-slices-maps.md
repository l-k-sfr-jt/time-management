# 05 — Slices & maps

## Goal

Get comfortable with Go's two everyday collection types before we start
returning lists of Groups/Activity Types from the database.

## Concepts

- Arrays (fixed size, rarely used directly) vs. **slices** (dynamically
  sized, what you actually use — like a resizable array/list).
- `append`, and the "sometimes it copies, sometimes it doesn't" mental
  model of slice growth (capacity vs. length).
- `nil` slices vs. empty slices — both usable, subtly different, matters
  for JSON output (`null` vs `[]`).
- **Maps**: `map[string]int`, checking for a missing key with the
  two-value form (`v, ok := m["key"]`), zero value on missing key.
- Iterating with `range` over slices and maps (and why map iteration
  order is intentionally randomized).

*(Task and checkpoint written when we start this lesson.)*

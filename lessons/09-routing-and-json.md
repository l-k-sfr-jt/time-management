# 09 — Routing & JSON

## Goal

Go from one hardcoded handler to real routing with path parameters, and
learn how Go structs turn into JSON (and back).

## Concepts

- Go 1.22+'s `http.ServeMux` supports method + path patterns directly:
  `mux.HandleFunc("PATCH /groups/{id}", handler)`, with `r.PathValue("id")`
  to read the wildcard — no router library needed (see our earlier
  chat about chi vs. stdlib).
- `encoding/json`: struct tags (`` `json:"name"` ``) controlling field
  names, `json.NewEncoder(w).Encode(v)` to write a response,
  `json.NewDecoder(r.Body).Decode(&v)` to read a request body.
- Pointers in JSON structs (`*string`) as the idiomatic way to represent
  "field not provided" vs. "field explicitly set to empty string" —
  you'll see this all over `internal/httpapi/dto.go` in the reference
  implementation.

*(Task and checkpoint written when we start this lesson.)*

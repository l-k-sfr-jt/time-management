# 08 — net/http fundamentals

## Goal

Build a real, if tiny, HTTP server: `GET /healthz` returning `200 OK`.
This is the foundation everything in Part 2 sits on.

## Concepts

- `http.Handler` — the interface at the center of Go's HTTP story: any
  type with `ServeHTTP(http.ResponseWriter, *http.Request)` can handle a
  request. (This is the "implicit interfaces" idea from lesson 04, put
  to use.)
- `http.HandlerFunc` — an adapter that lets a plain function satisfy
  `http.Handler`, so you rarely write `ServeHTTP` by hand.
- `http.ResponseWriter` (write status + body) and `*http.Request` (read
  method, headers, body, URL).
- `http.ListenAndServe` — the simplest way to start serving.
- Why no framework is needed for this part — it's genuinely this small
  in the standard library.

*(Task and checkpoint written when we start this lesson.)*

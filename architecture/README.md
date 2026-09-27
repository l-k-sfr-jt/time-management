# Architecture

Technical design docs for the time management app, written the way you'd
lay out a system-design interview: entities → data flows → API → deep
dives. Product-level functional/non-functional requirements live in
`../docs/requirements.md`; this folder is the "how we build it" layer on
top of those.

- **[data-model.md](./data-model.md)** — entities, Postgres schema,
  recurrence-materialization design, Clerk sync notes.
- **[data-flows.md](./data-flows.md)** — sequence diagrams for the flows
  that aren't obvious from the schema alone: recording an activity,
  creating/editing a recurring plan, computing a summary, sending a
  reminder, and syncing a user from Clerk.
- **[api-design.md](./api-design.md)** — the Go backend's API, in the
  requirements → entities → API → high-level design → deep dives format:
  endpoint table, component diagram, and the trade-offs behind the harder
  parts (recurrence materialization at scale, overlapping timers, summary
  query performance, notification delivery reliability).

Stack recap (full rationale in `docs/requirements.md` §9): Go backend,
Postgres via Neon, TanStack Start frontend, Clerk for auth.

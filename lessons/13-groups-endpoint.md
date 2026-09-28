# 13 — Groups endpoint, end to end

## Goal

Put lessons 08–12 together: a real `POST/GET/PATCH/DELETE /groups`
implementation backed by Postgres. This is the first "full vertical
slice" — the payoff lesson for Part 2 so far.

## Concepts

- No new language concepts — this is where everything so far gets
  combined: routing (09) + JSON (09) + config (10) + pgx (11) + sqlc (12).
- The FR-1.3 rule (a Group can't be deleted while Activity Types still
  reference it) as a small piece of real business logic sitting in a
  handler, not the database.
- Matching a spec: `architecture/api-design.md` §3.1 already defines the
  exact request/response shapes — the task is building to match a
  contract, like a real ticket would.

*(Task and checkpoint written when we start this lesson.)*

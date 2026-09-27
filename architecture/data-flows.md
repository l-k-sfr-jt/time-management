# Data Flows

Sequence diagrams for the flows that touch multiple components. Entities
referenced here are defined in [data-model.md](./data-model.md); endpoints
are defined in [api-design.md](./api-design.md).

## 1. Start / stop recording an activity

The most latency-sensitive flow (NFR-1: <200ms perceived) — it's a single
row write in both directions, no joins needed on the hot path.

```mermaid
sequenceDiagram
    actor U as User
    participant C as Client (TanStack Start)
    participant A as Go API
    participant DB as Postgres

    U->>C: Tap "Start" on Activity Type X
    C->>A: POST /time-logs/start {activityTypeId}
    A->>DB: INSERT time_logs (start_at=now(), end_at=NULL)
    DB-->>A: log row
    A-->>C: 201 {id, activityTypeId, startAt}
    C-->>U: show running timer (elapsed ticking client-side)

    Note over C,A: Multiple timers can be running at once (FR-4.4) —<br/>no check against other running logs.

    U->>C: Tap "Stop"
    C->>A: POST /time-logs/{id}/stop
    A->>DB: UPDATE time_logs SET end_at = now() WHERE id = $1 AND user_id = $2
    DB-->>A: updated row
    A-->>C: 200 {id, startAt, endAt, durationMinutes}
```

Retroactive logging (FR-4.5) is the same write, just with client-supplied
`startAt`/`endAt` via `POST /time-logs` instead of the start/stop pair.

## 2. Create a recurring plan (and materialize occurrences)

```mermaid
sequenceDiagram
    actor U as User
    participant C as Client
    participant A as Go API
    participant DB as Postgres

    U->>C: Define plan: Activity X, weekly Mon/Wed/Fri, 07:00, 45 min
    C->>A: POST /planned-series {activityTypeId, recurrence...}
    A->>DB: INSERT planned_series
    A->>A: generateOccurrences(series, from=today, to=today+60d)
    A->>DB: BULK INSERT planned_occurrences (one row per date the rule matches)
    DB-->>A: ok
    A-->>C: 201 {seriesId, occurrenceCount}
    C-->>U: show new entries on calendar
```

A daily background job extends every active series' materialized window
forward (see [api-design.md §6.1](./api-design.md#61-recurrence-materialization-at-scale))
so the client never has to interpret recurrence rules itself.

### 2a. Editing one occurrence ("this occurrence only" — the default)

```mermaid
sequenceDiagram
    actor U as User
    participant C as Client
    participant A as Go API
    participant DB as Postgres

    U->>C: Move Wed's occurrence to 08:00 (not the whole series)
    C->>A: PATCH /planned-occurrences/{id}?scope=this {startAt: ...}
    A->>DB: UPDATE planned_occurrences SET start_at=..., is_overridden=true WHERE id=$1
    A-->>C: 200 updated occurrence
    Note over A,DB: planned_series is untouched — future Mon/Fri occurrences unaffected.
```

"This and following" / "all" scopes are handled per
[data-model.md §3](./data-model.md#3-recurrence-materialization) — they
touch `planned_series`, not just one row, and trigger re-materialization.

## 3. Planned-vs-actual summary (week/month view)

```mermaid
sequenceDiagram
    actor U as User
    participant C as Client
    participant A as Go API
    participant DB as Postgres

    U->>C: Open "This Week" summary
    C->>A: GET /summary?period=week&date=2026-09-27
    A->>DB: SELECT group_id, activity_type_id, SUM(duration) FROM planned_occurrences WHERE user_id=$1 AND occurrence_date BETWEEN $2 AND $3 AND NOT is_canceled GROUP BY ...
    A->>DB: SELECT group_id, activity_type_id, SUM(EXTRACT(EPOCH FROM (end_at-start_at))) FROM time_logs WHERE user_id=$1 AND start_at BETWEEN $2 AND $3 GROUP BY ...
    A->>DB: SELECT planned_occurrences LEFT JOIN time_logs (via planned_occurrence_id) WHERE time_logs.id IS NULL -- "missed"
    A->>DB: SELECT time_logs WHERE planned_occurrence_id IS NULL -- "unplanned"
    A->>A: merge planned + actual by (group, activity), compute variance
    A-->>C: 200 {byGroup: [...], missed: [...], unplanned: [...]}
    C-->>U: render chart + variance list
```

Both aggregate queries run independently (they hit different tables) and
are merged in the API layer rather than with one large SQL join — keeps
each query simple and independently indexable
(see [api-design.md §6.2](./api-design.md#62-summary-query-performance)).

## 4. Reminder notification before a planned entry

```mermaid
sequenceDiagram
    participant W as Notification Worker (Go, ticks every minute)
    participant DB as Postgres
    participant N as Push/Email provider

    loop every 60s
        W->>DB: SELECT * FROM planned_occurrences<br/>WHERE notification_sent_at IS NULL<br/>AND start_at - lead_minutes <= now()<br/>AND start_at > now()<br/>AND NOT is_canceled
        DB-->>W: due occurrences
        W->>N: send notification(user, occurrence)
        N-->>W: ack
        W->>DB: UPDATE planned_occurrences SET notification_sent_at = now() WHERE id = $1
    end
```

`notification_sent_at` makes delivery idempotent — a worker crash mid-batch
just means some rows retry the next tick, never double-sent (barring a
crash between send and stamp, which is an accepted at-least-once edge
case — see [api-design.md §6.4](./api-design.md#64-notification-delivery-reliability)).

## 5. User sync from Clerk

```mermaid
sequenceDiagram
    actor U as User
    participant Clerk
    participant C as Client
    participant A as Go API
    participant DB as Postgres

    U->>Clerk: Sign up / log in
    Clerk-->>C: session JWT
    Clerk->>A: webhook: user.created {clerkUserId, email}
    A->>DB: UPSERT users (clerk_user_id, email)

    Note over C,A: Normal authenticated request, any endpoint above:
    C->>A: GET /groups  (Authorization: Bearer <Clerk JWT>)
    A->>A: verify JWT against Clerk JWKS
    A->>DB: SELECT id FROM users WHERE clerk_user_id = $1
    DB-->>A: internal user_id
    A->>DB: ... query scoped to user_id ...
    A-->>C: 200 response
```

If a request arrives before the webhook has landed (race on first login),
the API upserts the `users` row itself from the verified JWT claims instead
of failing — the webhook is an optimization, not the only path.

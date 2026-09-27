# Data Model

Target: Postgres (Neon). Backend: Go, using `pgx`/`sqlc`-style plain SQL
(no heavy ORM). All tables are scoped by `user_id`; there is no
cross-user sharing in v1.

## 1. Entity Overview

```
users
  └─< groups
        └─< activity_types
                ├─< planned_series ──< planned_occurrences ─── time_logs (optional link)
                └─< time_logs
```

- **users** — synced from Clerk (Clerk is the source of truth for
  credentials; we keep a local row per user for foreign keys and app data).
- **groups** — topics (Work, Sport, Relax, ...), user-defined.
- **activity_types** — named activities, each belongs to one group.
- **planned_series** — the *definition* of a plan: either a one-off plan or
  a recurring plan (holds the recurrence rule inline).
- **planned_occurrences** — concrete, materialized calendar instances
  generated from a series (this is what the planning/calendar UI and the
  summary queries actually read). Supports per-occurrence overrides and
  cancellation without mutating the whole series.
- **time_logs** — actual recorded activity (start/end), optionally linked to
  a `planned_occurrence`. Overlapping time ranges are allowed (no exclusion
  constraint).

Series + materialized occurrences is the standard pattern for recurring
calendar events: it makes "this occurrence only" edits/deletes (FR-3.5/3.6)
simple (mutate one `planned_occurrences` row) without needing to
re-interpret the recurrence rule on every read.

## 2. Tables

### users
Synced from Clerk via webhook (`user.created`/`user.updated`/`user.deleted`).

| Column | Type | Notes |
|---|---|---|
| id | uuid, PK | internal id |
| clerk_user_id | text, unique, not null | Clerk's `user.id` |
| email | text, not null | from Clerk, for display only |
| timezone | text, not null, default `'UTC'` | IANA tz name, e.g. `Europe/Bratislava`; used for day/week/month boundaries (NFR re: DST) |
| notification_lead_minutes | int, not null, default 10 | default reminder lead time (FR-3.9), overridable per series |
| created_at | timestamptz, not null, default now() | |
| updated_at | timestamptz, not null, default now() | |

### groups
| Column | Type | Notes |
|---|---|---|
| id | uuid, PK | |
| user_id | uuid, FK → users, not null | |
| name | text, not null | |
| color | text, nullable | hex color |
| icon | text, nullable | icon key |
| is_default | boolean, not null, default false | seeded defaults (Work/Sport/Relax/Study/Hobbies), user can still edit/delete (FR-1.5) |
| created_at / updated_at | timestamptz | |

- Unique index on `(user_id, lower(name))`.
- Deleting a group is blocked at the application layer while
  `activity_types` reference it (FR-1.3) — enforced with `ON DELETE
  RESTRICT` on the FK from `activity_types.group_id`.

### activity_types
| Column | Type | Notes |
|---|---|---|
| id | uuid, PK | |
| user_id | uuid, FK → users, not null | denormalized for query simplicity/RLS-style filtering |
| group_id | uuid, FK → groups, not null, `ON DELETE RESTRICT` | |
| name | text, not null | |
| description | text, nullable | |
| default_duration_minutes | int, nullable | prefill for planning/quick-start |
| color | text, nullable | overrides group color if set |
| archived_at | timestamptz, nullable | archiving instead of hard delete (FR-2.3); archived types are hidden from pickers but historical `time_logs`/`planned_occurrences` keep referencing them |
| created_at / updated_at | timestamptz | |

- Non-unique index on `(user_id, group_id, lower(name))` — duplicates are
  allowed but flagged as a warning in the UI (FR-2.5), not DB-enforced.

### planned_series
The definition of a plan — a single occurrence or a recurring rule.

| Column | Type | Notes |
|---|---|---|
| id | uuid, PK | |
| user_id | uuid, FK → users, not null | |
| activity_type_id | uuid, FK → activity_types, not null | |
| start_date | date, not null | first occurrence's date |
| start_time | time, not null | local time of day the activity starts |
| duration_minutes | int, not null | |
| recurrence_frequency | text, nullable | `null` = one-off; else `daily`, `weekly`, `monthly_by_date`, `monthly_by_weekday`, `custom_interval` |
| recurrence_interval | int, nullable | "every N" — days for `daily`/`custom_interval`, weeks for `weekly` |
| recurrence_weekdays | smallint[], nullable | for `weekly`, ISO weekdays 1–7 |
| recurrence_month_day | int, nullable | for `monthly_by_date` |
| recurrence_week_ordinal | int, nullable | for `monthly_by_weekday`, e.g. `3` = third |
| recurrence_weekday | smallint, nullable | for `monthly_by_weekday`, ISO weekday 1–7 |
| recurrence_end_type | text, nullable | `never` \| `on_date` \| `after_count` |
| recurrence_end_date | date, nullable | |
| recurrence_end_count | int, nullable | |
| notification_lead_minutes | int, nullable | overrides `users.notification_lead_minutes` when set (FR-3.9) |
| canceled_from_date | date, nullable | set when an "this and following" edit/delete truncates this series (see §3) |
| created_at / updated_at | timestamptz | |

### planned_occurrences
Materialized, concrete calendar instances generated from a series. This is
what planning/calendar views and summary queries read directly — no
recurrence-rule interpretation needed at read time.

| Column | Type | Notes |
|---|---|---|
| id | uuid, PK | |
| user_id | uuid, FK → users, not null | |
| series_id | uuid, FK → planned_series, not null, `ON DELETE CASCADE` | |
| activity_type_id | uuid, FK → activity_types, not null | copied from series at generation time; can differ if this single occurrence was edited (override) |
| occurrence_date | date, not null | |
| start_at | timestamptz, not null | derived from `occurrence_date` + `start_time` + user's timezone at generation time |
| end_at | timestamptz, not null | |
| is_overridden | boolean, not null, default false | true once this row was edited independently of its series ("this occurrence only") |
| is_canceled | boolean, not null, default false | true if this single occurrence was deleted ("this occurrence only") |
| notification_sent_at | timestamptz, nullable | dedupes reminder delivery |
| created_at / updated_at | timestamptz | |

- Unique index on `(series_id, occurrence_date)`.
- Index on `(user_id, start_at)` for calendar/summary range queries.
- A rolling window (e.g. next 60 days) is (re)generated by a backend job
  whenever a series is created/edited, and extended periodically — see §3.

### time_logs
| Column | Type | Notes |
|---|---|---|
| id | uuid, PK | |
| user_id | uuid, FK → users, not null | |
| activity_type_id | uuid, FK → activity_types, not null | |
| start_at | timestamptz, not null | |
| end_at | timestamptz, nullable | `null` while the timer is running (FR-4.2) |
| note | text, nullable | FR-4.8 |
| planned_occurrence_id | uuid, FK → planned_occurrences, nullable | link for "done as planned" (FR-4.6); nullable/unlinkable |
| created_at / updated_at | timestamptz | |

- Index on `(user_id, start_at)`.
- Partial index on `(user_id) WHERE end_at IS NULL` to quickly list
  currently-running timers (FR-4.4).
- **No** exclusion/overlap constraint — overlapping logs are allowed by
  requirement.

## 3. Recurrence Materialization

- On create/edit of a `planned_series`, the backend generates
  `planned_occurrences` rows for a rolling window (e.g. today → +60 days,
  or up to `recurrence_end_date`/`recurrence_end_count` if sooner).
- A scheduled job extends the window forward periodically (e.g. daily) so
  the calendar never runs out of materialized occurrences.
- Edit/delete scope handling (FR-3.5/3.6), default = "this occurrence only":
  - **This occurrence only**: mutate/cancel the single
    `planned_occurrences` row (`is_overridden`/`is_canceled = true`);
    `planned_series` untouched.
  - **This and following occurrences**: set `planned_series.canceled_from_date`
    on the original series to the edited occurrence's date, delete its
    future (unmodified) occurrences from that date on, and create a **new**
    `planned_series` starting at that date with the updated fields —
    regenerating occurrences for it.
  - **All occurrences**: update `planned_series` in place and regenerate
    all of its non-overridden future occurrences.

## 4. Notifications (FR-3.9)

A scheduled worker polls `planned_occurrences` for rows where
`start_at - notification_lead_minutes` (falling back to the user's default)
has just passed and `notification_sent_at IS NULL`, sends the notification,
and stamps `notification_sent_at`.

## 5. Clerk Integration Notes

- Clerk issues a session JWT; the Go backend verifies it per-request (Clerk
  Go SDK / JWKS) and resolves `users.id` from `clerk_user_id` (claim `sub`).
- A Clerk webhook (`user.created`) upserts the local `users` row so
  foreign keys exist before the user's first authenticated request.

## 6. Future: AI Chat Extensibility

Not designed yet, but the schema stays compatible with adding, later:

```
conversations (id, user_id, title, created_at, ...)
messages (id, conversation_id, role, content, created_at, ...)
```

Both would be `user_id`-scoped like every other table here, so no existing
table needs to change shape to support it.

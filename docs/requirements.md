# Time Management App — Functional Requirements

## 1. Purpose

A personal application to plan, record, and review daily activities. It lets the
user log what they actually did throughout the day, organize activities into
topic groups (work, sport, relax, study, hobbies, ...), plan activities ahead
of time (including recurring ones), and compare planned vs. actual time usage
through weekly/monthly summaries.

## 2. Actors

- **User** — single-user app (no multi-tenant/collaboration requirements at
  this stage). All requirements below assume one authenticated user per
  account.

## 3. Core Concepts / Domain Model

| Concept | Description |
|---|---|
| **Group (Topic)** | A category activities belong to, e.g. Work, Sport, Relax, Study, Hobbies. User-defined, not a fixed list. |
| **Activity Type** | A named kind of activity (e.g. "Running", "Reading", "Team standup"), always assigned to exactly one Group. Reusable template for both planning and recording. |
| **Planned Entry** | A future/scheduled occurrence of an Activity Type, either one-off or generated from a recurrence rule. |
| **Recurrence Rule** | Defines how a Planned Entry repeats (daily, weekly on certain days, monthly, custom interval, with optional end date/count). |
| **Time Log (Recorded Activity)** | An actual, timestamped occurrence of an Activity Type: start time, end time (or duration), optionally linked to the Planned Entry it fulfills. |
| **Summary** | Aggregated view (by day/week/month, and by group/activity) comparing planned time vs. logged time. |

## 4. Functional Requirements

### 4.1 Group / Topic Management
- FR-1.1: User can create a new Group with a name and (optional) color/icon.
- FR-1.2: User can rename, recolor, or delete an existing Group.
- FR-1.3: Deleting a Group that still has Activity Types must be blocked, or
  require reassigning/archiving its Activity Types first (avoid orphaned data).
- FR-1.4: User can view a list of all Groups.
- FR-1.5: System ships with a small set of default Groups (Work, Sport, Relax,
  Study, Hobbies) on first use, which the user can edit or delete like any
  other Group.

### 4.2 Activity Type Management
- FR-2.1: User can create an Activity Type with: name (required), Group
  (required), optional description/notes, optional default duration, optional
  color override.
- FR-2.2: User can edit or archive/delete an Activity Type.
- FR-2.3: Archiving (rather than hard-deleting) an Activity Type preserves
  historical Time Logs and Planned Entries that reference it.
- FR-2.4: User can browse/search Activity Types, filterable by Group.
- FR-2.5: Activity Type names should be unique within a Group (warn on
  duplicates, don't hard-block).

### 4.3 Planning
- FR-3.1: User can create a Planned Entry for a specific Activity Type on a
  specific date, with a start time and either an end time or a duration.
- FR-3.2: User can edit or delete a Planned Entry.
- FR-3.3: User can view planned entries in a calendar/agenda view (day, week,
  month).
- FR-3.4: User can create a **recurring** Planned Entry by attaching a
  Recurrence Rule:
  - Frequency: daily, weekly (with selectable weekdays), monthly (by date or
    by weekday-in-month), custom interval (every N days/weeks).
  - Optional end condition: never, on a specific date, or after N
    occurrences.
- FR-3.5: Editing a recurring series lets the user choose the scope: "this
  occurrence only", "this and following occurrences", or "all occurrences".
- FR-3.6: Deleting a recurring series has the same scope options as editing.
- FR-3.7: The system materializes/generates concrete occurrences from a
  Recurrence Rule far enough ahead to populate calendar views (e.g. rolling
  window, generated on demand rather than infinitely stored).
- FR-3.8: User can see, for a given day/week, the list of planned activities
  and their total planned time, optionally grouped by Group.
- FR-3.9: The system sends a notification a configurable short time before a
  Planned Entry's start time (default e.g. 5–10 minutes prior); the user can
  adjust or disable this per entry or globally.

### 4.4 Recording Activities (Tracking Screen)
- FR-4.1: A dedicated "Record" screen lets the user start tracking an
  activity in real time: pick (or quick-search) an Activity Type, tap Start,
  and the app timestamps the start.
- FR-4.2: While an activity is running, the screen shows a live elapsed-time
  display and a Stop control.
- FR-4.3: Stopping the timer creates a Time Log with start time, end time, and
  computed duration.
- FR-4.4: Multiple activities can be tracked concurrently — starting a new
  timer does not stop any other currently running timer. The recording
  screen lists all currently running timers together, each with its own
  Stop control.
- FR-4.5: User can manually add/edit a Time Log after the fact (retroactive
  logging) with a custom start/end time or duration, for cases where they
  forgot to start the timer.
- FR-4.6: If a Time Log's date/activity matches an existing Planned Entry for
  that day, the app suggests linking the log to that plan (so it counts as
  "done as planned"); user can also link/unlink manually.
- FR-4.7: The recording screen surfaces quick-start shortcuts for the user's
  most-used or most-recent Activity Types, and shows today's Planned Entries
  as one-tap "start this planned activity" actions.
- FR-4.8: User can add a short free-text note to a Time Log (e.g. what
  specifically was done).
- FR-4.9: User can delete a Time Log.

### 4.5 Daily / Weekly / Monthly Summary
- FR-5.1: User can view a summary for a selected day, week, or month showing
  total time logged, broken down by Group and by Activity Type.
- FR-5.2: Summary includes a **planned vs. actual** comparison: for each
  Group/Activity Type, planned time vs. logged time, with the
  difference/variance highlighted (e.g. under/over, or not started).
- FR-5.3: Summary lists Planned Entries that had no matching Time Log
  ("missed"/"not done") and Time Logs with no matching plan
  ("unplanned"/"extra").
- FR-5.4: User can visualize the breakdown (e.g. bar chart or pie chart by
  Group) for the selected period.
- FR-5.5: User can navigate between periods (previous/next day, week, month)
  and jump to "today/this week/this month".
- FR-5.6: Summary numbers are filterable by Group (e.g. "show only Work").

### 4.6 Cross-cutting
- FR-6.1: All timestamps are stored with enough precision to compute
  durations and handle the user's local timezone correctly, including DST
  transitions.
- FR-6.2: User authentication/login is required (single-user account) since
  data is synced via a backend and accessed from multiple devices; handled by
  Clerk (see Section 9, Technical Stack).
- FR-6.3: Data export (CSV/JSON) is out of scope for this iteration — see
  Section 7.

## 5. Non-Functional Requirements

- NFR-1: The recording screen (start/stop) must respond instantly (<200ms
  perceived) since it's used many times a day.
- NFR-2: Works well on mobile (primary use case is logging activities
  on-the-go) as well as desktop/web for planning and reviewing summaries.
- NFR-3: Data persists reliably; no data loss if the app is closed while a
  timer is running (timer state must survive app restarts).
- NFR-4: Summary/report queries over a month of data should load in <1s.

## 6. Decisions

These were open questions, now resolved and reflected in the requirements
above:

| Question | Decision |
|---|---|
| Platform for v1 | **Web app** (responsive, usable on phone browser and desktop). |
| Overlapping Time Logs | **Allowed.** The user may have more than one activity logged/running with overlapping time ranges (e.g. "commute" + "podcast"). FR-4.4 is updated accordingly (see below). |
| Time-budget goals vs. scheduled plans | **Scheduled plans only for v1** (specific date/time entries, with recurrence). Standalone weekly/monthly time budgets not tied to a slot are out of scope for now. |
| Multi-device sync | **Cloud sync with login.** The app requires an account; data is stored server-side and accessible from any device. |
| Reminders/notifications | **Included in v1** — basic notification shortly before a planned entry's start time. |
| Offline support | **Not required.** App assumes online connectivity; no offline queue/conflict handling needed. |
| Data export | **Deferred**, not required for v1. |
| Default recurring-edit scope | **"This occurrence only"** is the default/first offered option when editing or deleting an occurrence of a recurring plan (FR-3.5/3.6). |

All affected requirements above (FR-3.9, FR-4.4, FR-6.2, FR-6.3) already
reflect these decisions.

## 7. Out of Scope (for this iteration)

- Multi-user collaboration or sharing activities/plans with others.
- Third-party calendar sync (Google Calendar, Outlook) — may be considered
  later.
- Offline recording/sync.
- Data export (CSV/JSON).
- Standalone time-budget goals independent of scheduled plans.
- Gamification (streaks, badges).

## 8. Suggested Implementation Plan (high level)

1. **Data model & backend**: Groups, Activity Types, Recurrence Rules,
   Planned Entries (materialized occurrences), Time Logs (supporting
   overlapping ranges), user accounts/auth. Basic CRUD APIs.
2. **Recording screen**: start/stop timer(s) — supports multiple concurrent
   running activities — quick-start list, manual entry.
3. **Planning screen**: calendar/agenda view, create/edit one-off and
   recurring plans (default edit scope: this occurrence only).
4. **Group/Activity management screens**: CRUD for Groups and Activity
   Types.
5. **Reminders**: scheduled notification before a Planned Entry's start
   time.
6. **Summary/reports screen**: day/week/month aggregation, planned-vs-actual
   comparison, charts.
7. **Polish**: timezone/DST handling, auth hardening.

## 9. Technical Stack

| Layer | Choice | Notes |
|---|---|---|
| Backend | **Go** (REST API) | Single service for v1; router TBD (e.g. chi/std `net/http`). |
| Database | **Postgres (Neon)** | Serverless Postgres, already provisioned in this environment. |
| DB access | **pgx** + plain SQL / `sqlc` | Avoids heavy ORM; keeps recurrence-rule and time-range queries explicit and reviewable. |
| Frontend | **TanStack Start** (React) | File-based routing, SSR, type-safe data loading via TanStack Router/Query. |
| Auth | **Clerk** | Managed auth + authorization (sessions, user management); Go backend verifies Clerk session JWTs, frontend uses Clerk's React SDK. |
| Notifications (FR-3.9) | Web push / in-app, delivery mechanism TBD | Needs a scheduler (e.g. a lightweight cron/worker in the Go service) that checks upcoming Planned Entries. |
| Future | AI chat feature | Not designed yet; data model keeps all tables `user_id`-scoped so a `conversations`/`messages` area can be added later without reshaping existing tables. |

See `docs/data-model.md` for the concrete schema.

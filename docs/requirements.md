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

### 4.4 Recording Activities (Tracking Screen)
- FR-4.1: A dedicated "Record" screen lets the user start tracking an
  activity in real time: pick (or quick-search) an Activity Type, tap Start,
  and the app timestamps the start.
- FR-4.2: While an activity is running, the screen shows a live elapsed-time
  display and a Stop control.
- FR-4.3: Stopping the timer creates a Time Log with start time, end time, and
  computed duration.
- FR-4.4: Only one activity can be actively running at a time (starting a new
  one prompts to stop/replace the current one, or the app auto-stops it).
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
- FR-6.2: Basic data export (e.g. CSV/JSON) of Time Logs for a date range.
- FR-6.3: User authentication/login (single user, but app should not be
  wide open) — exact mechanism to be decided in technical design.

## 5. Non-Functional Requirements

- NFR-1: The recording screen (start/stop) must respond instantly (<200ms
  perceived) since it's used many times a day.
- NFR-2: Works well on mobile (primary use case is logging activities
  on-the-go) as well as desktop/web for planning and reviewing summaries.
- NFR-3: Data persists reliably; no data loss if the app is closed while a
  timer is running (timer state must survive app restarts).
- NFR-4: Summary/report queries over a month of data should load in <1s.

## 6. Out of Scope (for this iteration)

- Multi-user collaboration or sharing activities/plans with others.
- Third-party calendar sync (Google Calendar, Outlook) — may be considered
  later.
- Notifications/reminders for planned activities — candidate for a later
  iteration.
- Gamification (streaks, badges).

## 7. Open Questions

1. Should overlapping Time Logs be allowed (e.g. logging two things at once,
   like "commute" + "listening to podcast"), or strictly one active/logged
   activity at a time?
2. Should Planned Entries support a target duration goal per week/month (e.g.
   "Sport: 3h/week") independent of specific calendar slots, in addition to
   scheduled entries?
3. What platform(s) first — mobile app, web app, or both from the start?
4. Any need for reminders/notifications before a planned activity starts?

## 8. Suggested Implementation Plan (high level)

1. **Data model & backend**: Groups, Activity Types, Recurrence Rules,
   Planned Entries (materialized occurrences), Time Logs. Basic CRUD APIs.
2. **Recording screen**: start/stop timer, quick-start list, manual entry.
3. **Planning screen**: calendar/agenda view, create/edit one-off and
   recurring plans.
4. **Group/Activity management screens**: CRUD for Groups and Activity
   Types.
5. **Summary/reports screen**: day/week/month aggregation, planned-vs-actual
   comparison, charts.
6. **Polish**: export, timezone/DST handling, offline-safe timer state.

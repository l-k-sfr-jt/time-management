-- Initial schema for the time management app.
-- See docs/data-model.md for the full design rationale.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id                          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id               text NOT NULL UNIQUE,
    email                       text NOT NULL,
    timezone                    text NOT NULL DEFAULT 'UTC',
    notification_lead_minutes  int NOT NULL DEFAULT 10,
    created_at                  timestamptz NOT NULL DEFAULT now(),
    updated_at                  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE groups (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        text NOT NULL,
    color       text,
    icon        text,
    is_default  boolean NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX groups_user_name_uidx ON groups (user_id, lower(name));

CREATE TABLE activity_types (
    id                          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id                    uuid NOT NULL REFERENCES groups(id) ON DELETE RESTRICT,
    name                        text NOT NULL,
    description                 text,
    default_duration_minutes   int,
    color                       text,
    archived_at                 timestamptz,
    created_at                  timestamptz NOT NULL DEFAULT now(),
    updated_at                  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX activity_types_user_group_name_idx
    ON activity_types (user_id, group_id, lower(name));

CREATE TABLE planned_series (
    id                          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type_id            uuid NOT NULL REFERENCES activity_types(id) ON DELETE RESTRICT,
    start_date                  date NOT NULL,
    start_time                  time NOT NULL,
    duration_minutes            int NOT NULL,
    recurrence_frequency        text
        CHECK (recurrence_frequency IN ('daily', 'weekly', 'monthly_by_date', 'monthly_by_weekday', 'custom_interval')),
    recurrence_interval          int,
    recurrence_weekdays          smallint[],
    recurrence_month_day         int,
    recurrence_week_ordinal      int,
    recurrence_weekday           smallint,
    recurrence_end_type          text
        CHECK (recurrence_end_type IN ('never', 'on_date', 'after_count')),
    recurrence_end_date          date,
    recurrence_end_count         int,
    notification_lead_minutes    int,
    canceled_from_date           date,
    created_at                   timestamptz NOT NULL DEFAULT now(),
    updated_at                   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX planned_series_user_idx ON planned_series (user_id);

CREATE TABLE planned_occurrences (
    id                      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id               uuid NOT NULL REFERENCES planned_series(id) ON DELETE CASCADE,
    activity_type_id        uuid NOT NULL REFERENCES activity_types(id) ON DELETE RESTRICT,
    occurrence_date         date NOT NULL,
    start_at                timestamptz NOT NULL,
    end_at                  timestamptz NOT NULL,
    is_overridden           boolean NOT NULL DEFAULT false,
    is_canceled             boolean NOT NULL DEFAULT false,
    notification_sent_at    timestamptz,
    created_at              timestamptz NOT NULL DEFAULT now(),
    updated_at              timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX planned_occurrences_series_date_uidx
    ON planned_occurrences (series_id, occurrence_date);
CREATE INDEX planned_occurrences_user_start_idx
    ON planned_occurrences (user_id, start_at);

CREATE TABLE time_logs (
    id                      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type_id        uuid NOT NULL REFERENCES activity_types(id) ON DELETE RESTRICT,
    start_at                timestamptz NOT NULL,
    end_at                  timestamptz,
    note                    text,
    planned_occurrence_id   uuid REFERENCES planned_occurrences(id) ON DELETE SET NULL,
    created_at              timestamptz NOT NULL DEFAULT now(),
    updated_at              timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX time_logs_user_start_idx ON time_logs (user_id, start_at);
CREATE INDEX time_logs_running_idx ON time_logs (user_id) WHERE end_at IS NULL;

-- name: CreateActivityType :one
INSERT INTO activity_types (user_id, group_id, name, description, default_duration_minutes, color)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListActivityTypes :many
SELECT * FROM activity_types
WHERE user_id = $1
    AND (sqlc.narg(group_id)::uuid IS NULL OR group_id = sqlc.narg(group_id))
    AND (sqlc.arg(include_archived)::bool OR archived_at IS NULL)
ORDER BY created_at;

-- name: GetActivityType :one
SELECT * FROM activity_types WHERE id = $1 AND user_id = $2;

-- name: UpdateActivityType :one
UPDATE activity_types
SET
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    default_duration_minutes = COALESCE(sqlc.narg(default_duration_minutes), default_duration_minutes),
    color = COALESCE(sqlc.narg(color), color),
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: ArchiveActivityType :one
UPDATE activity_types
SET archived_at = now(), updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

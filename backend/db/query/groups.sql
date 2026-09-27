-- name: CreateGroup :one
INSERT INTO groups (user_id, name, color, icon)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListGroups :many
SELECT * FROM groups
WHERE user_id = $1
ORDER BY created_at;

-- name: GetGroup :one
SELECT * FROM groups WHERE id = $1 AND user_id = $2;

-- name: UpdateGroup :one
UPDATE groups
SET
    name = COALESCE(sqlc.narg(name), name),
    color = COALESCE(sqlc.narg(color), color),
    icon = COALESCE(sqlc.narg(icon), icon),
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteGroup :execrows
DELETE FROM groups WHERE id = $1 AND user_id = $2;

-- name: CountActivityTypesInGroup :one
SELECT count(*) FROM activity_types WHERE group_id = $1;

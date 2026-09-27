-- name: UpsertUser :one
INSERT INTO users (clerk_user_id, email)
VALUES ($1, $2)
ON CONFLICT (clerk_user_id) DO UPDATE SET
    email = EXCLUDED.email,
    updated_at = now()
RETURNING *;

-- name: GetUserByClerkID :one
SELECT * FROM users WHERE clerk_user_id = $1;

-- name: UpdateUserSettings :one
UPDATE users
SET
    timezone = COALESCE(sqlc.narg(timezone), timezone),
    notification_lead_minutes = COALESCE(sqlc.narg(notification_lead_minutes), notification_lead_minutes),
    updated_at = now()
WHERE id = $1
RETURNING *;

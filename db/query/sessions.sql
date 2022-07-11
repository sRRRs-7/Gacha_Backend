-- name: CreateSession :one
INSERT INTO sessions (
    user_name,
    user_agent,
    client_ip,
    is_blocked,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;
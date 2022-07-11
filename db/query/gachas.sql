-- name: CreateGacha :one
INSERT INTO gachas (
    account_id, item_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetGacha :one
SELECT * FROM gachas
WHERE id = $1
LIMIT 1;

-- name: ListGachas :many
SELECT * FROM gachas
ORDER BY id ASC
LIMIT $1
OFFSET $2;
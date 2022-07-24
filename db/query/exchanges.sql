-- name: CreateExchange :one
INSERT INTO exchanges (
    from_account_id, to_account_id, item_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetExchange :one
SELECT * FROM exchanges
WHERE id = $1 LIMIT 1;

-- name: ListExchangeFromAccount :many
SELECT * FROM exchanges
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListExchangeToAccount :many
SELECT * FROM exchanges
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
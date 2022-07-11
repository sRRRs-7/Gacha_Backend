-- name: CreateItem :one
INSERT INTO items (
    item_name, rating, item_url, category_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = $1 LIMIT 1;

-- name: ListItemByCategoryId :many
SELECT * FROM items
WHERE category_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListItemsByCategoryId :many
SELECT * FROM items
ORDER BY category_id ASC
LIMIT $1
OFFSET $2;

-- name: ListItemsById :many
SELECT * FROM items
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListItemsByItemName :many
SELECT * FROM items
ORDER BY item_name ASC
LIMIT $1
OFFSET $2;

-- name: ListItemsByRating :many
SELECT * FROM items
ORDER BY rating DESC
LIMIT $1
OFFSET $2;

-- name: UpdateItem :one
UPDATE items
SET item_name = $2, rating = $3, item_url = $4, category_id = $5
where id = $1
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1;
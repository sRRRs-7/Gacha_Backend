-- name: CreateCategory :one
INSERT INTO categories (
    category
) VALUES (
    $1
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE category = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY category ASC
LIMIT $1
OFFSET $2;
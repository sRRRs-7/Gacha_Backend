-- name: CreateGallery :one
INSERT INTO galleries (
    owner_id, item_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetGallery :one
SELECT * FROM galleries
WHERE id = $1 LIMIT 1;

-- name: ListGalleriesById :many
SELECT * FROM galleries
WHERE owner_id = $1
ORDER BY id ASC
LIMIT $2
OFFSET $3;

-- name: ListGalleriesByItemId :many
SELECT * FROM galleries
WHERE item_id = $1
ORDER BY item_id ASC
LIMIT $2
OFFSET $3;
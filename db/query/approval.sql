-- name: CreateApproval :one
INSERT INTO approval (
    from_account_id,
    from_item_id,
    from_A_approval,
    to_account_id,
    to_item_id,
    to_A_approval
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetApproval :one
SELECT * FROM approval
WHERE id = $1 LIMIT 1;

-- name: ListApproval :many
SELECT * FROM approval
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateApprovalRequest :one
UPDATE approval
SET from_A_approval = $2
where id = $1
RETURNING *;

-- name: UpdateApprovalResponse :one
UPDATE approval
SET to_A_approval = $2
where id = $1
RETURNING *;

-- name: DeleteApproval :exec
DELETE FROM approval
WHERE id = $1;
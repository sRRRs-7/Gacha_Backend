// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: categories.sql

package db

import (
	"context"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (
    category
) VALUES (
    $1
) RETURNING id, category, created_at
`

func (q *Queries) CreateCategory(ctx context.Context, category string) (Category, error) {
	row := q.db.QueryRowContext(ctx, createCategory, category)
	var i Category
	err := row.Scan(&i.ID, &i.Category, &i.CreatedAt)
	return i, err
}

const getCategory = `-- name: GetCategory :one
SELECT id, category, created_at FROM categories
WHERE category = $1 LIMIT 1
`

func (q *Queries) GetCategory(ctx context.Context, category string) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategory, category)
	var i Category
	err := row.Scan(&i.ID, &i.Category, &i.CreatedAt)
	return i, err
}

const listCategories = `-- name: ListCategories :many
SELECT id, category, created_at FROM categories
ORDER BY category ASC
LIMIT $1
OFFSET $2
`

type ListCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, listCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Category, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

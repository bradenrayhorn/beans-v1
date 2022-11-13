// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: month_categories.sql

package db

import (
	"context"

	"github.com/jackc/pgtype"
)

const createMonthCategory = `-- name: CreateMonthCategory :exec
INSERT INTO month_categories (
  id, month_id, category_id, amount
) VALUES ($1, $2, $3, $4)
`

type CreateMonthCategoryParams struct {
	ID         string
	MonthID    string
	CategoryID string
	Amount     pgtype.Numeric
}

func (q *Queries) CreateMonthCategory(ctx context.Context, arg CreateMonthCategoryParams) error {
	_, err := q.db.Exec(ctx, createMonthCategory,
		arg.ID,
		arg.MonthID,
		arg.CategoryID,
		arg.Amount,
	)
	return err
}

const getMonthCategoriesForMonth = `-- name: GetMonthCategoriesForMonth :many
SELECT id, month_id, category_id, amount, created_at FROM month_categories WHERE month_id = $1
`

func (q *Queries) GetMonthCategoriesForMonth(ctx context.Context, monthID string) ([]MonthCategory, error) {
	rows, err := q.db.Query(ctx, getMonthCategoriesForMonth, monthID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthCategory
	for rows.Next() {
		var i MonthCategory
		if err := rows.Scan(
			&i.ID,
			&i.MonthID,
			&i.CategoryID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonthCategoryByMonthAndCategory = `-- name: GetMonthCategoryByMonthAndCategory :one
SELECT id, month_id, category_id, amount, created_at FROM month_categories WHERE month_id = $1 and category_id = $2
`

type GetMonthCategoryByMonthAndCategoryParams struct {
	MonthID    string
	CategoryID string
}

func (q *Queries) GetMonthCategoryByMonthAndCategory(ctx context.Context, arg GetMonthCategoryByMonthAndCategoryParams) (MonthCategory, error) {
	row := q.db.QueryRow(ctx, getMonthCategoryByMonthAndCategory, arg.MonthID, arg.CategoryID)
	var i MonthCategory
	err := row.Scan(
		&i.ID,
		&i.MonthID,
		&i.CategoryID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const updateMonthCategoryAmount = `-- name: UpdateMonthCategoryAmount :exec
UPDATE month_categories SET amount = $1 WHERE id = $2
`

type UpdateMonthCategoryAmountParams struct {
	Amount pgtype.Numeric
	ID     string
}

func (q *Queries) UpdateMonthCategoryAmount(ctx context.Context, arg UpdateMonthCategoryAmountParams) error {
	_, err := q.db.Exec(ctx, updateMonthCategoryAmount, arg.Amount, arg.ID)
	return err
}

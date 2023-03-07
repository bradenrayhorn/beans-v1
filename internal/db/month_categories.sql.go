// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: month_categories.sql

package db

import (
	"context"
	"time"

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

const getAssignedInMonth = `-- name: GetAssignedInMonth :one
SELECT sum(month_categories.amount)::numeric as amount
  FROM month_categories
  JOIN months m on m.id = month_categories.month_id
    AND m.id = $1
`

func (q *Queries) GetAssignedInMonth(ctx context.Context, id string) (pgtype.Numeric, error) {
	row := q.db.QueryRow(ctx, getAssignedInMonth, id)
	var amount pgtype.Numeric
	err := row.Scan(&amount)
	return amount, err
}

const getMonthCategoriesForMonth = `-- name: GetMonthCategoriesForMonth :many
SELECT month_categories.id, month_categories.month_id, month_categories.category_id, month_categories.amount, month_categories.created_at, sum(t.amount)::numeric as activity
  FROM month_categories
  LEFT JOIN transactions t on t.category_id = month_categories.category_id
    AND t.date >= $1 AND t.date <= $2
  WHERE month_id = $3
  GROUP BY (
    month_categories.id,
    month_categories.month_id,
    month_categories.category_id,
    month_categories.amount
  )
`

type GetMonthCategoriesForMonthParams struct {
	FromDate time.Time
	ToDate   time.Time
	MonthID  string
}

type GetMonthCategoriesForMonthRow struct {
	ID         string
	MonthID    string
	CategoryID string
	Amount     pgtype.Numeric
	CreatedAt  time.Time
	Activity   pgtype.Numeric
}

func (q *Queries) GetMonthCategoriesForMonth(ctx context.Context, arg GetMonthCategoriesForMonthParams) ([]GetMonthCategoriesForMonthRow, error) {
	rows, err := q.db.Query(ctx, getMonthCategoriesForMonth, arg.FromDate, arg.ToDate, arg.MonthID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMonthCategoriesForMonthRow
	for rows.Next() {
		var i GetMonthCategoriesForMonthRow
		if err := rows.Scan(
			&i.ID,
			&i.MonthID,
			&i.CategoryID,
			&i.Amount,
			&i.CreatedAt,
			&i.Activity,
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

const getPastMonthCategoriesAvailable = `-- name: GetPastMonthCategoriesAvailable :many
SELECT
    categories.id,
    sum(mc.amount)::numeric as assigned
  FROM categories
  JOIN month_categories mc on mc.category_id = categories.id
  JOIN months m on m.id = mc.month_id
    AND m.budget_id = $1
    AND m.date < $2
  GROUP BY (
    categories.id
  )
`

type GetPastMonthCategoriesAvailableParams struct {
	BudgetID   string
	BeforeDate time.Time
}

type GetPastMonthCategoriesAvailableRow struct {
	ID       string
	Assigned pgtype.Numeric
}

func (q *Queries) GetPastMonthCategoriesAvailable(ctx context.Context, arg GetPastMonthCategoriesAvailableParams) ([]GetPastMonthCategoriesAvailableRow, error) {
	rows, err := q.db.Query(ctx, getPastMonthCategoriesAvailable, arg.BudgetID, arg.BeforeDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPastMonthCategoriesAvailableRow
	for rows.Next() {
		var i GetPastMonthCategoriesAvailableRow
		if err := rows.Scan(&i.ID, &i.Assigned); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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

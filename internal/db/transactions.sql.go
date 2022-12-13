// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: transactions.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
)

const createTransaction = `-- name: CreateTransaction :exec
INSERT INTO transactions (
  id, account_id, category_id, date, amount, notes
) VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateTransactionParams struct {
	ID         string
	AccountID  string
	CategoryID sql.NullString
	Date       time.Time
	Amount     pgtype.Numeric
	Notes      sql.NullString
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) error {
	_, err := q.db.Exec(ctx, createTransaction,
		arg.ID,
		arg.AccountID,
		arg.CategoryID,
		arg.Date,
		arg.Amount,
		arg.Notes,
	)
	return err
}

const getTransactionsForBudget = `-- name: GetTransactionsForBudget :many
SELECT transactions.id, transactions.account_id, transactions.payee_id, transactions.category_id, transactions.date, transactions.amount, transactions.notes, transactions.created_at, accounts.name as account_name, categories.name as category_name from transactions
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = $1
LEFT JOIN categories
  ON categories.id = transactions.category_id
`

type GetTransactionsForBudgetRow struct {
	ID           string
	AccountID    string
	PayeeID      sql.NullString
	CategoryID   sql.NullString
	Date         time.Time
	Amount       pgtype.Numeric
	Notes        sql.NullString
	CreatedAt    time.Time
	AccountName  string
	CategoryName sql.NullString
}

func (q *Queries) GetTransactionsForBudget(ctx context.Context, budgetID string) ([]GetTransactionsForBudgetRow, error) {
	rows, err := q.db.Query(ctx, getTransactionsForBudget, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransactionsForBudgetRow
	for rows.Next() {
		var i GetTransactionsForBudgetRow
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.PayeeID,
			&i.CategoryID,
			&i.Date,
			&i.Amount,
			&i.Notes,
			&i.CreatedAt,
			&i.AccountName,
			&i.CategoryName,
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

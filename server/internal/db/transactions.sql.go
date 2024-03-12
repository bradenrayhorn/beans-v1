// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: transactions.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTransaction = `-- name: CreateTransaction :exec
INSERT INTO transactions (
  id, account_id, payee_id, category_id, date, amount, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreateTransactionParams struct {
	ID         string
	AccountID  string
	PayeeID    pgtype.Text
	CategoryID pgtype.Text
	Date       pgtype.Date
	Amount     pgtype.Numeric
	Notes      pgtype.Text
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) error {
	_, err := q.db.Exec(ctx, createTransaction,
		arg.ID,
		arg.AccountID,
		arg.PayeeID,
		arg.CategoryID,
		arg.Date,
		arg.Amount,
		arg.Notes,
	)
	return err
}

const deleteTransactions = `-- name: DeleteTransactions :exec
DELETE FROM transactions
  USING accounts
  WHERE
    accounts.id = transactions.account_id
    AND accounts.budget_id=$1
    AND transactions.id = ANY($2::varchar[])
`

type DeleteTransactionsParams struct {
	BudgetID string
	Ids      []string
}

func (q *Queries) DeleteTransactions(ctx context.Context, arg DeleteTransactionsParams) error {
	_, err := q.db.Exec(ctx, deleteTransactions, arg.BudgetID, arg.Ids)
	return err
}

const getActivityByCategory = `-- name: GetActivityByCategory :many
SELECT categories.id, sum(transactions.amount)::numeric as activity
  FROM transactions
  JOIN categories
    ON transactions.category_id = categories.id
  JOIN accounts
    ON accounts.id = transactions.account_id
    AND accounts.budget_id = $1
  WHERE
    (transactions.date >= $2 OR NOT $3)
    AND (transactions.date <= $4 OR NOT $5)
  GROUP BY (
    categories.id
  )
`

type GetActivityByCategoryParams struct {
	BudgetID       string
	FromDate       pgtype.Date
	FilterFromDate interface{}
	ToDate         pgtype.Date
	FilterToDate   interface{}
}

type GetActivityByCategoryRow struct {
	ID       string
	Activity pgtype.Numeric
}

func (q *Queries) GetActivityByCategory(ctx context.Context, arg GetActivityByCategoryParams) ([]GetActivityByCategoryRow, error) {
	rows, err := q.db.Query(ctx, getActivityByCategory,
		arg.BudgetID,
		arg.FromDate,
		arg.FilterFromDate,
		arg.ToDate,
		arg.FilterToDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetActivityByCategoryRow
	for rows.Next() {
		var i GetActivityByCategoryRow
		if err := rows.Scan(&i.ID, &i.Activity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIncomeBetween = `-- name: GetIncomeBetween :one
SELECT sum(transactions.amount)::numeric
FROM transactions
JOIN categories
  ON categories.id = transactions.category_id
JOIN category_groups
  ON category_groups.id = categories.group_id
  AND category_groups.is_income = true
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = $1
WHERE
  transactions.date <= $2
  AND transactions.date >= $3
`

type GetIncomeBetweenParams struct {
	BudgetID  string
	EndDate   pgtype.Date
	BeginDate pgtype.Date
}

func (q *Queries) GetIncomeBetween(ctx context.Context, arg GetIncomeBetweenParams) (pgtype.Numeric, error) {
	row := q.db.QueryRow(ctx, getIncomeBetween, arg.BudgetID, arg.EndDate, arg.BeginDate)
	var column_1 pgtype.Numeric
	err := row.Scan(&column_1)
	return column_1, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT transactions.id, transactions.account_id, transactions.payee_id, transactions.category_id, transactions.date, transactions.amount, transactions.notes, transactions.created_at
  FROM transactions
  JOIN accounts
    ON accounts.id = transactions.account_id
    AND accounts.budget_id = $1
  WHERE transactions.id = $2
`

type GetTransactionParams struct {
	BudgetID string
	ID       string
}

func (q *Queries) GetTransaction(ctx context.Context, arg GetTransactionParams) (Transaction, error) {
	row := q.db.QueryRow(ctx, getTransaction, arg.BudgetID, arg.ID)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.PayeeID,
		&i.CategoryID,
		&i.Date,
		&i.Amount,
		&i.Notes,
		&i.CreatedAt,
	)
	return i, err
}

const getTransactionsForBudget = `-- name: GetTransactionsForBudget :many
SELECT
  transactions.id, transactions.account_id, transactions.payee_id, transactions.category_id, transactions.date, transactions.amount, transactions.notes, transactions.created_at,
  accounts.name as account_name,
  categories.name as category_name,
  payees.name as payee_name,
  accounts.off_budget as account_off_budget
FROM transactions
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = $1
LEFT JOIN categories
  ON categories.id = transactions.category_id
LEFT JOIN payees
  ON payees.id = transactions.payee_id
ORDER BY date desc
`

type GetTransactionsForBudgetRow struct {
	Transaction      Transaction
	AccountName      string
	CategoryName     pgtype.Text
	PayeeName        pgtype.Text
	AccountOffBudget bool
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
			&i.Transaction.ID,
			&i.Transaction.AccountID,
			&i.Transaction.PayeeID,
			&i.Transaction.CategoryID,
			&i.Transaction.Date,
			&i.Transaction.Amount,
			&i.Transaction.Notes,
			&i.Transaction.CreatedAt,
			&i.AccountName,
			&i.CategoryName,
			&i.PayeeName,
			&i.AccountOffBudget,
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

const updateTransaction = `-- name: UpdateTransaction :exec
UPDATE transactions
  SET account_id=$1, category_id=$2, payee_id=$3, date=$4, amount=$5, notes=$6
  WHERE id=$7
`

type UpdateTransactionParams struct {
	AccountID  string
	CategoryID pgtype.Text
	PayeeID    pgtype.Text
	Date       pgtype.Date
	Amount     pgtype.Numeric
	Notes      pgtype.Text
	ID         string
}

func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) error {
	_, err := q.db.Exec(ctx, updateTransaction,
		arg.AccountID,
		arg.CategoryID,
		arg.PayeeID,
		arg.Date,
		arg.Amount,
		arg.Notes,
		arg.ID,
	)
	return err
}

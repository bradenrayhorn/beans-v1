// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
)

type Account struct {
	ID        string
	Name      string
	BudgetID  string
	CreatedAt time.Time
}

type Budget struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type BudgetsUser struct {
	BudgetID string
	UserID   string
}

type Category struct {
	ID        string
	Name      string
	BudgetID  string
	GroupID   string
	CreatedAt time.Time
}

type CategoryGroup struct {
	ID        string
	Name      string
	IsIncome  bool
	BudgetID  string
	CreatedAt time.Time
}

type Month struct {
	ID        string
	BudgetID  string
	Date      time.Time
	Carryover pgtype.Numeric
	CreatedAt time.Time
}

type MonthCategory struct {
	ID         string
	MonthID    string
	CategoryID string
	Amount     pgtype.Numeric
	CreatedAt  time.Time
}

type Payee struct {
	ID        string
	Name      string
	BudgetID  string
	CreatedAt time.Time
}

type Transaction struct {
	ID         string
	AccountID  string
	PayeeID    sql.NullString
	CategoryID sql.NullString
	Date       time.Time
	Amount     pgtype.Numeric
	Notes      sql.NullString
	CreatedAt  time.Time
}

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"time"
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

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

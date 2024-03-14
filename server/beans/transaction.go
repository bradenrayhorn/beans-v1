package beans

import (
	"context"
)

type Transaction struct {
	ID ID

	AccountID  ID
	CategoryID ID
	PayeeID    ID

	Amount Amount
	Date   Date
	Notes  TransactionNotes

	TransferID ID
}

type TransactionWithRelations struct {
	Transaction

	Variant TransactionVariant

	Account  RelatedAccount
	Category Optional[RelatedCategory]
	Payee    Optional[RelatedPayee]

	TransferAccount Optional[RelatedAccount]
}

type TransactionNotes struct{ NullString }

func NewTransactionNotes(string string) TransactionNotes {
	return TransactionNotes{NullString: NewNullString(string)}
}

type TransactionVariant string

const (
	TransactionStandard  TransactionVariant = "standard"
	TransactionOffBudget TransactionVariant = "off_budget"
	TransactionTransfer  TransactionVariant = "transfer"
)

type TransactionContract interface {
	// Creates a transaction.
	Create(ctx context.Context, auth *BudgetAuthContext, params TransactionCreateParams) (ID, error)

	// Gets all transactions for budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]TransactionWithRelations, error)

	// Edits a transaction.
	Update(ctx context.Context, auth *BudgetAuthContext, params TransactionUpdateParams) error

	// Deletes transactions.
	Delete(ctx context.Context, auth *BudgetAuthContext, transactionIDs []ID) error

	// Gets a transaction details.
	Get(ctx context.Context, auth *BudgetAuthContext, id ID) (TransactionWithRelations, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transactions []Transaction) error

	Update(ctx context.Context, transactions []Transaction) error

	Delete(ctx context.Context, budgetID ID, transactionIDs []ID) error

	GetForBudget(ctx context.Context, budgetID ID) ([]TransactionWithRelations, error)

	GetWithRelations(ctx context.Context, budgetID ID, id ID) (TransactionWithRelations, error)

	// Get transaction.
	Get(ctx context.Context, budgetID ID, id ID) (Transaction, error)

	// Gets sum of all income transactions between the dates.
	GetIncomeBetween(ctx context.Context, budgetID ID, begin Date, end Date) (Amount, error)

	// Gets sum of transactions grouped by category between the dates.
	GetActivityByCategory(ctx context.Context, budgetID ID, from Date, to Date) (map[ID]Amount, error)
}

type TransactionParams struct {
	AccountID  ID
	CategoryID ID
	PayeeID    ID
	Amount     Amount
	Date       Date
	Notes      TransactionNotes
}

type TransactionCreateParams struct {
	TransferAccountID ID
	TransactionParams
}

func (t TransactionCreateParams) ValidateAll() error {
	if err := t.TransactionParams.ValidateAll(); err != nil {
		return err
	}

	if !t.TransferAccountID.Empty() && (!t.PayeeID.Empty() || !t.CategoryID.Empty()) {
		return NewError(EINVALID, "cannot set a payee or category on transfer")
	}

	return nil
}

type TransactionUpdateParams struct {
	ID ID
	TransactionParams
}

func (t TransactionUpdateParams) ValidateAll() error {
	if err := t.TransactionParams.ValidateAll(); err != nil {
		return err
	}

	return ValidateFields(Field("Transaction ID", Required(t.ID)))
}

func (t TransactionParams) ValidateAll() error {
	return ValidateFields(
		Field("Account ID", Required(t.AccountID)),
		Field("Amount", Required(&t.Amount), MaxPrecision(t.Amount)),
		Field("Date", Required(t.Date)),
		Field("Notes", Max(t.Notes, 255, "characters")),
	)
}

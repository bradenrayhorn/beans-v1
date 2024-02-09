package httpadapter

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
)

func mapAll[T any, K any](objs []T, mapper func(T) K) []K {
	var models []K
	for _, m := range objs {
		mapped := mapper(m)

		models = append(models, mapped)
	}

	return models
}

func mapBudget(t response.Budget) beans.Budget {
	return beans.Budget{
		ID:   t.ID,
		Name: beans.Name(t.Name),
	}
}

func mapCategory(t response.Category) beans.Category {
	return beans.Category{
		ID:   t.ID,
		Name: beans.Name(t.Name),
	}
}

func mapCategoryGroup(t response.CategoryGroup) beans.CategoryGroup {
	return beans.CategoryGroup{
		ID:       t.ID,
		Name:     beans.Name(t.Name),
		IsIncome: t.IsIncome,
	}
}

func mapAccount(t response.Account) beans.Account {
	return beans.Account{ID: t.ID, Name: beans.Name(t.Name)}
}

func mapListAccount(t response.ListAccount) beans.AccountWithBalance {
	return beans.AccountWithBalance{
		Account: beans.Account{ID: t.ID, Name: beans.Name(t.Name)},
		Balance: t.Balance,
	}
}

func mapTransactionWithRelations(t response.Transaction) beans.TransactionWithRelations {
	transaction := beans.TransactionWithRelations{
		Transaction: beans.Transaction{
			ID:        t.ID,
			AccountID: t.Account.ID,
			Amount:    t.Amount,
			Date:      t.Date,
			Notes:     t.Notes,
		},
		Account: beans.RelatedAccount{
			ID:   t.Account.ID,
			Name: t.Account.Name,
		},
	}

	if t.Category != nil {
		transaction.Category = beans.OptionalWrap(beans.RelatedCategory{ID: t.Category.ID, Name: t.Category.Name})
	}
	if t.Payee != nil {
		transaction.Payee = beans.OptionalWrap(beans.RelatedPayee{ID: t.Payee.ID, Name: t.Payee.Name})
	}

	return transaction
}

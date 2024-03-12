package datasource

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTransaction(t *testing.T, ds beans.DataSource) {
	factory := testutils.NewFactory(t, ds)

	transactionRepository := ds.TransactionRepository()
	ctx := context.Background()

	t.Run("can create", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})
		payee := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		err := transactionRepository.Create(
			ctx,
			beans.Transaction{
				ID:         beans.NewID(),
				AccountID:  account.ID,
				CategoryID: category.ID,
				PayeeID:    payee.ID,
				Amount:     beans.NewAmount(5, 0),
				Date:       beans.NewDate(time.Now()),
				Notes:      beans.NewTransactionNotes("notes"),
			},
		)
		require.Nil(t, err)
	})

	t.Run("can create with empty optional fields", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})

		transaction1 := beans.Transaction{
			ID:        beans.NewID(),
			AccountID: account.ID,
			Amount:    beans.NewAmount(5, 0),
			Date:      testutils.NewDate(t, "2022-08-28"),
		}
		require.Nil(t, transactionRepository.Create(ctx, transaction1))

		transactions, err := transactionRepository.GetForBudget(ctx, budget.ID)
		require.Nil(t, err)
		assert.Len(t, transactions, 1)

		assert.True(t, transactions[0].CategoryID.Empty())
		assert.True(t, transactions[0].PayeeID.Empty())
		assert.True(t, transactions[0].Notes.Empty())
	})

	t.Run("cannot get nonexistant", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		_, err := transactionRepository.Get(ctx, budget.ID, beans.NewID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("cannot get for other budget", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()
		transaction := factory.Transaction(budget2.ID, beans.Transaction{})

		_, err := transactionRepository.Get(ctx, budget.ID, transaction.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})
		payee := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		transaction := beans.Transaction{
			ID:         beans.NewID(),
			AccountID:  account.ID,
			CategoryID: category.ID,
			PayeeID:    payee.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
		}
		require.Nil(t, transactionRepository.Create(ctx, transaction))

		res, err := transactionRepository.Get(ctx, budget.ID, transaction.ID)
		require.Nil(t, err)

		assert.Equal(t, transaction, res)
	})

	t.Run("can update", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account1 := factory.Account(beans.Account{BudgetID: budget.ID})
		account2 := factory.Account(beans.Account{BudgetID: budget.ID})
		payee1 := factory.Payee(beans.Payee{BudgetID: budget.ID})
		payee2 := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category1 := factory.Category(beans.Category{BudgetID: budget.ID})
		category2 := factory.Category(beans.Category{BudgetID: budget.ID})

		transaction := beans.Transaction{
			ID:         beans.NewID(),
			AccountID:  account1.ID,
			CategoryID: category1.ID,
			PayeeID:    payee1.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
		}
		require.Nil(t, transactionRepository.Create(ctx, transaction))

		transaction.AccountID = account2.ID
		transaction.CategoryID = category2.ID
		transaction.PayeeID = payee2.ID
		transaction.Amount = beans.NewAmount(6, 0)
		transaction.Date = testutils.NewDate(t, "2022-08-30")
		transaction.Notes = beans.NewTransactionNotes("notes 5")

		require.Nil(t, transactionRepository.Update(ctx, transaction))

		res, err := transactionRepository.Get(ctx, budget.ID, transaction.ID)
		require.Nil(t, err)

		assert.True(t, reflect.DeepEqual(transaction, res))
	})

	t.Run("can delete", func(t *testing.T) {
		budget1, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()

		transaction1 := factory.Transaction(budget1.ID, beans.Transaction{})
		transaction2 := factory.Transaction(budget1.ID, beans.Transaction{})
		transaction3 := factory.Transaction(budget1.ID, beans.Transaction{})
		transaction4 := factory.Transaction(budget2.ID, beans.Transaction{})

		err := transactionRepository.Delete(ctx, budget1.ID, []beans.ID{transaction1.ID, transaction2.ID, transaction4.ID})
		require.Nil(t, err)

		// transaction1 and transaction2 should be deleted, they are passed in and part of budget 1.
		// transaction3 should not be deleted, it is not passed in.
		// transaction4 should not be deleted, it is passed in but not part of budget 1.
		_, err = transactionRepository.Get(ctx, budget1.ID, transaction1.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		_, err = transactionRepository.Get(ctx, budget1.ID, transaction2.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		_, err = transactionRepository.Get(ctx, budget1.ID, transaction3.ID)
		assert.Nil(t, err)

		_, err = transactionRepository.Get(ctx, budget2.ID, transaction4.ID)
		assert.Nil(t, err)
	})

	t.Run("get all for budget", func(t *testing.T) {

		t.Run("can get all", func(t *testing.T) {
			budget1, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget1.ID})
			payee := factory.Payee(beans.Payee{BudgetID: budget1.ID})
			category := factory.Category(beans.Category{BudgetID: budget1.ID})

			transaction1 := factory.Transaction(budget1.ID, beans.Transaction{
				AccountID:  account.ID,
				PayeeID:    payee.ID,
				CategoryID: category.ID,
			})
			// this transaction should not be included
			factory.Transaction(budget2.ID, beans.Transaction{})

			transactions, err := transactionRepository.GetForBudget(ctx, budget1.ID)
			require.Nil(t, err)
			assert.Len(t, transactions, 1)

			assert.Equal(t, beans.TransactionWithRelations{
				Transaction: transaction1,
				Variant:     beans.TransactionStandard,
				Account:     beans.RelatedAccount{ID: account.ID, Name: account.Name},
				Category:    beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name}),
				Payee:       beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name}),
			}, transactions[0])
		})

		t.Run("maps off budget variant", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID, OffBudget: true})

			transaction := factory.Transaction(budget.ID, beans.Transaction{AccountID: account.ID})

			res, err := transactionRepository.GetForBudget(ctx, budget.ID)
			require.Nil(t, err)
			require.Equal(t, 1, len(res))

			assert.Equal(t, beans.TransactionWithRelations{
				Transaction: transaction,
				Variant:     beans.TransactionOffBudget,
				Account:     beans.RelatedAccount{ID: account.ID, Name: account.Name},
				Category:    beans.Optional[beans.RelatedCategory]{},
				Payee:       beans.Optional[beans.RelatedPayee]{},
			}, res[0])
		})
	})

	t.Run("can get activity by category", func(t *testing.T) {

		t.Run("groups and sums", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			category1 := factory.Category(beans.Category{BudgetID: budget.ID})
			category2 := factory.Category(beans.Category{BudgetID: budget.ID})

			// setup 3 transactions - two in category1 and one in category2
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				CategoryID: category1.ID,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(2, 0),
				CategoryID: category1.ID,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(1, 0),
				CategoryID: category2.ID,
			})

			// make sure they are grouped and summed properly
			res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, beans.Date{})
			require.NoError(t, err)

			assert.Equal(t, 2, len(res))
			assert.Equal(t, beans.NewAmount(5, 0), res[category1.ID])
			assert.Equal(t, beans.NewAmount(1, 0), res[category2.ID])
		})

		t.Run("filters by date", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			// setup 4 transactions with varying dates
			category := factory.Category(beans.Category{BudgetID: budget.ID})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2099-09-01"),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(2, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-09-01"),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(1, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-08-31"),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(8, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "1900-08-31"),
			})

			// try to filter the transactions

			t.Run("filters by only from date", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, testutils.NewDate(t, "2022-09-01"), beans.Date{})
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(5, 0), res[category.ID])
			})

			t.Run("filters by only to date", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, testutils.NewDate(t, "2022-08-31"))
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(9, 0), res[category.ID])
			})

			t.Run("filters by both dates", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-09-30"))
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(3, 0), res[category.ID])
			})

			t.Run("applies no date filter", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, beans.Date{})
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(14, 0), res[category.ID])
			})
		})

		t.Run("filters by budget", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			factory.Transaction(budget2.ID, beans.Transaction{
				Amount: beans.NewAmount(1, 0),
				Date:   testutils.NewDate(t, "2022-09-01"),
			})

			res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, beans.Date{})
			require.NoError(t, err)

			assert.Equal(t, 0, len(res))
		})
	})

	t.Run("can get income", func(t *testing.T) {

		t.Run("can get", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			incomeGroup := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget.ID, IsIncome: true})
			incomeCategory := factory.Category(beans.Category{BudgetID: budget.ID, GroupID: incomeGroup.ID})
			otherCategory := factory.Category(beans.Category{BudgetID: budget.ID})

			budget2, _ := factory.MakeBudgetAndUser()
			budget2IncomeGroup := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget2.ID, IsIncome: true})
			budget2IncomeCategory := factory.Category(beans.Category{BudgetID: budget2.ID, GroupID: budget2IncomeGroup.ID})

			// Earned $1 in September
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(1, 0),
				Date:       testutils.NewDate(t, "2022-09-01"),
				CategoryID: incomeCategory.ID,
			})
			// Earned $2 in August
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(2, 0),
				Date:       testutils.NewDate(t, "2022-08-31"),
				CategoryID: incomeCategory.ID,
			})
			// Earned $3 in August
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				Date:       testutils.NewDate(t, "2022-08-01"),
				CategoryID: incomeCategory.ID,
			})
			// Earned $3 in July
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				Date:       testutils.NewDate(t, "2022-07-31"),
				CategoryID: incomeCategory.ID,
			})

			// Spent $99 in August
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(99, 0),
				Date:       testutils.NewDate(t, "2022-08-31"),
				CategoryID: otherCategory.ID,
			})
			// Spent $99 in July
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(99, 0),
				Date:       testutils.NewDate(t, "2022-07-29"),
				CategoryID: otherCategory.ID,
			})

			// Budget 2, earned $99 in August
			factory.Transaction(budget2.ID, beans.Transaction{
				Amount:     beans.NewAmount(99, 0),
				Date:       testutils.NewDate(t, "2022-08-15"),
				CategoryID: budget2IncomeCategory.ID,
			})

			amount, err := transactionRepository.GetIncomeBetween(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))
			require.Nil(t, err)

			// August earnings for budget 1 = $5
			require.Equal(t, beans.NewAmount(5, 0), amount)
		})

		t.Run("can get with no income", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			amount, err := transactionRepository.GetIncomeBetween(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))

			require.Nil(t, err)

			require.Equal(t, beans.NewAmount(0, 0), amount)
		})
	})
}

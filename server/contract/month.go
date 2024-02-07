package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type monthContract struct {
	contract
}

func (c *monthContract) GetOrCreate(ctx context.Context, auth *beans.BudgetAuthContext, date beans.MonthDate) (*beans.Month, []*beans.MonthCategory, beans.Amount, error) {
	month, err := c.createMonth(ctx, auth, date)
	if err != nil {
		return nil, nil, beans.NewEmptyAmount(), err
	}

	pastMonth, err := c.ds().MonthRepository().GetOrCreate(
		ctx,
		nil,
		auth.BudgetID(),
		month.Date.Previous(),
	)
	if err != nil {
		return nil, nil, beans.NewEmptyAmount(), err
	}

	categories, err := c.ds().MonthCategoryRepository().GetForMonth(ctx, month)
	if err != nil {
		return nil, nil, beans.NewEmptyAmount(), err
	}

	income, err := c.ds().TransactionRepository().GetIncomeBetween(ctx, auth.BudgetID(), month.Date.FirstDay(), month.Date.LastDay())
	if err != nil {
		return nil, nil, beans.NewEmptyAmount(), err
	}

	assignedInMonth, err := c.ds().MonthCategoryRepository().GetAssignedInMonth(ctx, month.ID)
	if err != nil {
		return nil, nil, beans.NewEmptyAmount(), err
	}

	available, err := beans.Arithmetic.Add(
		income,
		pastMonth.Carryover,
		beans.Arithmetic.Negate(month.Carryover),
		beans.Arithmetic.Negate(assignedInMonth),
	)
	if err != nil {
		return nil, nil, beans.NewEmptyAmount(), err
	}

	month.CarriedOver = pastMonth.Carryover
	month.Income = income
	month.Assigned = assignedInMonth

	return month, categories, available, nil
}

func (c *monthContract) createMonth(ctx context.Context, auth *beans.BudgetAuthContext, date beans.MonthDate) (*beans.Month, error) {
	return beans.ExecTx(ctx, c.ds().TxManager(), func(tx beans.Tx) (*beans.Month, error) {
		month, err := c.ds().MonthRepository().GetOrCreate(ctx, tx, auth.BudgetID(), date)
		if err != nil {
			return nil, err
		}

		categories, err := c.ds().CategoryRepository().GetForBudget(ctx, auth.BudgetID())
		if err != nil {
			return nil, err
		}

		for _, category := range categories {
			if _, err := c.ds().MonthCategoryRepository().GetOrCreate(ctx, tx, month.ID, category.ID); err != nil {
				return nil, err
			}
		}

		return month, nil
	})
}
func (c *monthContract) Update(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID, carryover beans.Amount) error {
	if err := beans.ValidateFields(
		beans.Field("Carryover", beans.Required(&carryover), beans.Positive(carryover)),
	); err != nil {
		return err
	}

	month, err := c.getAndVerifyMonth(ctx, auth, monthID)
	if err != nil {
		return err
	}

	month.Carryover = carryover

	return c.ds().MonthRepository().Update(ctx, month)
}

func (c *monthContract) SetCategoryAmount(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID, categoryID beans.ID, amount beans.Amount) error {
	if err := beans.ValidateFields(
		beans.Field("Amount", beans.NonZero(amount), beans.Positive(amount)),
	); err != nil {
		return err
	}

	_, err := c.getAndVerifyMonth(ctx, auth, monthID)
	if err != nil {
		return err
	}

	monthCategory, err := c.ds().MonthCategoryRepository().GetOrCreate(ctx, nil, monthID, categoryID)
	if err != nil {
		return err
	}

	return c.ds().MonthCategoryRepository().UpdateAmount(ctx, monthCategory.ID, amount)
}

func (c *monthContract) getAndVerifyMonth(ctx context.Context, auth *beans.BudgetAuthContext, monthID beans.ID) (*beans.Month, error) {
	month, err := c.ds().MonthRepository().Get(ctx, monthID)
	if err != nil {
		return nil, err
	}

	if err = c.verifyMonth(ctx, auth, month); err != nil {
		return nil, err
	}

	return month, nil
}

func (c *monthContract) verifyMonth(ctx context.Context, auth *beans.BudgetAuthContext, month *beans.Month) error {
	if month.BudgetID != auth.BudgetID() {
		return beans.NewError(beans.EFORBIDDEN, "No access to month")
	}

	return nil
}

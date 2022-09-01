// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	beans "github.com/bradenrayhorn/beans/beans"

	mock "github.com/stretchr/testify/mock"
)

// TransactionRepository is an autogenerated mock type for the TransactionRepository type
type TransactionRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, transaction
func (_m *TransactionRepository) Create(ctx context.Context, transaction *beans.Transaction) error {
	ret := _m.Called(ctx, transaction)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *beans.Transaction) error); ok {
		r0 = rf(ctx, transaction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetForBudget provides a mock function with given fields: ctx, budgetID
func (_m *TransactionRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Transaction, error) {
	ret := _m.Called(ctx, budgetID)

	var r0 []*beans.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, beans.ID) []*beans.Transaction); ok {
		r0 = rf(ctx, budgetID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*beans.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, beans.ID) error); ok {
		r1 = rf(ctx, budgetID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTransactionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionRepository creates a new instance of TransactionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionRepository(t mockConstructorTestingTNewTransactionRepository) *TransactionRepository {
	mock := &TransactionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

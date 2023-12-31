// Code generated by mockery v2.37.1. DO NOT EDIT.

package storage

import (
	domain "go-ticketos/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockOrderAdapter is an autogenerated mock type for the OrderAdapter type
type MockOrderAdapter struct {
	mock.Mock
}

// ToDomain provides a mock function with given fields: schema
func (_m *MockOrderAdapter) ToDomain(schema orderSchema) (*domain.Order, error) {
	ret := _m.Called(schema)

	var r0 *domain.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(orderSchema) (*domain.Order, error)); ok {
		return rf(schema)
	}
	if rf, ok := ret.Get(0).(func(orderSchema) *domain.Order); ok {
		r0 = rf(schema)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(orderSchema) error); ok {
		r1 = rf(schema)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToSchema provides a mock function with given fields: order
func (_m *MockOrderAdapter) ToSchema(order domain.Order) orderSchema {
	ret := _m.Called(order)

	var r0 orderSchema
	if rf, ok := ret.Get(0).(func(domain.Order) orderSchema); ok {
		r0 = rf(order)
	} else {
		r0 = ret.Get(0).(orderSchema)
	}

	return r0
}

// NewMockOrderAdapter creates a new instance of MockOrderAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOrderAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOrderAdapter {
	mock := &MockOrderAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

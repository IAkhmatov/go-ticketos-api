// Code generated by mockery v2.37.1. DO NOT EDIT.

package storage

import (
	domain "go-ticketos/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockTicketCategoryAdapter is an autogenerated mock type for the TicketCategoryAdapter type
type MockTicketCategoryAdapter struct {
	mock.Mock
}

// ToDomain provides a mock function with given fields: schema
func (_m *MockTicketCategoryAdapter) ToDomain(schema ticketCategorySchema) domain.TicketCategory {
	ret := _m.Called(schema)

	var r0 domain.TicketCategory
	if rf, ok := ret.Get(0).(func(ticketCategorySchema) domain.TicketCategory); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(domain.TicketCategory)
	}

	return r0
}

// ToSchema provides a mock function with given fields: tc
func (_m *MockTicketCategoryAdapter) ToSchema(tc domain.TicketCategory) ticketCategorySchema {
	ret := _m.Called(tc)

	var r0 ticketCategorySchema
	if rf, ok := ret.Get(0).(func(domain.TicketCategory) ticketCategorySchema); ok {
		r0 = rf(tc)
	} else {
		r0 = ret.Get(0).(ticketCategorySchema)
	}

	return r0
}

// NewMockTicketCategoryAdapter creates a new instance of MockTicketCategoryAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTicketCategoryAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTicketCategoryAdapter {
	mock := &MockTicketCategoryAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
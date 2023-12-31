// Code generated by mockery v2.37.1. DO NOT EDIT.

package storage

import (
	domain "go-ticketos/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockEventAdapter is an autogenerated mock type for the EventAdapter type
type MockEventAdapter struct {
	mock.Mock
}

// ToSchema provides a mock function with given fields: event
func (_m *MockEventAdapter) ToSchema(event domain.Event) eventSchema {
	ret := _m.Called(event)

	var r0 eventSchema
	if rf, ok := ret.Get(0).(func(domain.Event) eventSchema); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Get(0).(eventSchema)
	}

	return r0
}

// NewMockEventAdapter creates a new instance of MockEventAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventAdapter {
	mock := &MockEventAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

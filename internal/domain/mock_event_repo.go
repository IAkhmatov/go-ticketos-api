// Code generated by mockery v2.37.1. DO NOT EDIT.

package domain

import mock "github.com/stretchr/testify/mock"

// MockEventRepo is an autogenerated mock type for the EventRepo type
type MockEventRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: event
func (_m *MockEventRepo) Create(event Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockEventRepo creates a new instance of MockEventRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventRepo {
	mock := &MockEventRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

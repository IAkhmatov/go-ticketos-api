// Code generated by mockery v2.37.1. DO NOT EDIT.

package domain

import mock "github.com/stretchr/testify/mock"

// MockPayMasterWebHookUseCase is an autogenerated mock type for the PayMasterWebHookUseCase type
type MockPayMasterWebHookUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: props
func (_m *MockPayMasterWebHookUseCase) Execute(props PayMasterWebHookUseCaseProps) error {
	ret := _m.Called(props)

	var r0 error
	if rf, ok := ret.Get(0).(func(PayMasterWebHookUseCaseProps) error); ok {
		r0 = rf(props)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockPayMasterWebHookUseCase creates a new instance of MockPayMasterWebHookUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPayMasterWebHookUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPayMasterWebHookUseCase {
	mock := &MockPayMasterWebHookUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

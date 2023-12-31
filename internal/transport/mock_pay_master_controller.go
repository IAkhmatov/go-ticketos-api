// Code generated by mockery v2.37.1. DO NOT EDIT.

package transport

import mock "github.com/stretchr/testify/mock"

// MockPayMasterController is an autogenerated mock type for the PayMasterController type
type MockPayMasterController[context interface{}] struct {
	mock.Mock
}

// WebHook provides a mock function with given fields: c
func (_m *MockPayMasterController[context]) WebHook(c *context) error {
	ret := _m.Called(c)

	var r0 error
	if rf, ok := ret.Get(0).(func(*context) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockPayMasterController creates a new instance of MockPayMasterController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPayMasterController[context interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPayMasterController[context] {
	mock := &MockPayMasterController[context]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

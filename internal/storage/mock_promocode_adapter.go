// Code generated by mockery v2.37.1. DO NOT EDIT.

package storage

import (
	domain "go-ticketos/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockPromocodeAdapter is an autogenerated mock type for the PromocodeAdapter type
type MockPromocodeAdapter struct {
	mock.Mock
}

// ToDomain provides a mock function with given fields: promocode
func (_m *MockPromocodeAdapter) ToDomain(promocode promocodeSchema) domain.Promocode {
	ret := _m.Called(promocode)

	var r0 domain.Promocode
	if rf, ok := ret.Get(0).(func(promocodeSchema) domain.Promocode); ok {
		r0 = rf(promocode)
	} else {
		r0 = ret.Get(0).(domain.Promocode)
	}

	return r0
}

// ToSchema provides a mock function with given fields: promocode
func (_m *MockPromocodeAdapter) ToSchema(promocode domain.Promocode) promocodeSchema {
	ret := _m.Called(promocode)

	var r0 promocodeSchema
	if rf, ok := ret.Get(0).(func(domain.Promocode) promocodeSchema); ok {
		r0 = rf(promocode)
	} else {
		r0 = ret.Get(0).(promocodeSchema)
	}

	return r0
}

// NewMockPromocodeAdapter creates a new instance of MockPromocodeAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPromocodeAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPromocodeAdapter {
	mock := &MockPromocodeAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

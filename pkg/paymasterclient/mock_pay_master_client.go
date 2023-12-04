// Code generated by mockery v2.37.1. DO NOT EDIT.

package paymasterclient

import mock "github.com/stretchr/testify/mock"

// MockPayMasterClient is an autogenerated mock type for the PayMasterClient type
type MockPayMasterClient struct {
	mock.Mock
}

// CreateInvoice provides a mock function with given fields: props
func (_m *MockPayMasterClient) CreateInvoice(props CreateInvoiceRequestDTO) (*CreateInvoiceResponseDTO, error) {
	ret := _m.Called(props)

	var r0 *CreateInvoiceResponseDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(CreateInvoiceRequestDTO) (*CreateInvoiceResponseDTO, error)); ok {
		return rf(props)
	}
	if rf, ok := ret.Get(0).(func(CreateInvoiceRequestDTO) *CreateInvoiceResponseDTO); ok {
		r0 = rf(props)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CreateInvoiceResponseDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(CreateInvoiceRequestDTO) error); ok {
		r1 = rf(props)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockPayMasterClient creates a new instance of MockPayMasterClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPayMasterClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPayMasterClient {
	mock := &MockPayMasterClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
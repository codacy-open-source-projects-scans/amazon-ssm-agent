// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// ICredentialRefresher is an autogenerated mock type for the ICredentialRefresher type
type ICredentialRefresher struct {
	mock.Mock
}

// GetCredentialsReadyChan provides a mock function with given fields:
func (_m *ICredentialRefresher) GetCredentialsReadyChan() chan struct{} {
	ret := _m.Called()

	var r0 chan struct{}
	if rf, ok := ret.Get(0).(func() chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan struct{})
		}
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *ICredentialRefresher) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *ICredentialRefresher) Stop() {
	_m.Called()
}

// NewICredentialRefresher creates a new instance of ICredentialRefresher. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewICredentialRefresher(t testing.TB) *ICredentialRefresher {
	mock := &ICredentialRefresher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
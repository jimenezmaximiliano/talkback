// Code generated by mockery v2.32.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// LogError is an autogenerated mock type for the LogError type
type LogError struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, err
func (_m *LogError) Execute(ctx context.Context, err error) {
	_m.Called(ctx, err)
}

// NewLogError creates a new instance of LogError. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLogError(t interface {
	mock.TestingT
	Cleanup(func())
}) *LogError {
	mock := &LogError{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

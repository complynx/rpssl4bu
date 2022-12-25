// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// RandomProvider is an autogenerated mock type for the RandomProvider type
type RandomProvider struct {
	mock.Mock
}

type RandomProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *RandomProvider) EXPECT() *RandomProvider_Expecter {
	return &RandomProvider_Expecter{mock: &_m.Mock}
}

// Rand provides a mock function with given fields: ctx
func (_m *RandomProvider) Rand(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RandomProvider_Rand_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rand'
type RandomProvider_Rand_Call struct {
	*mock.Call
}

// Rand is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RandomProvider_Expecter) Rand(ctx interface{}) *RandomProvider_Rand_Call {
	return &RandomProvider_Rand_Call{Call: _e.mock.On("Rand", ctx)}
}

func (_c *RandomProvider_Rand_Call) Run(run func(ctx context.Context)) *RandomProvider_Rand_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RandomProvider_Rand_Call) Return(_a0 int, _a1 error) *RandomProvider_Rand_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewRandomProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewRandomProvider creates a new instance of RandomProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRandomProvider(t mockConstructorTestingTNewRandomProvider) *RandomProvider {
	mock := &RandomProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
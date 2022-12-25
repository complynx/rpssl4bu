// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/complynx/rpssl4bu/backend/pkg/types"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

type Storage_Expecter struct {
	mock *mock.Mock
}

func (_m *Storage) EXPECT() *Storage_Expecter {
	return &Storage_Expecter{mock: &_m.Mock}
}

// ClearScores provides a mock function with given fields:
func (_m *Storage) ClearScores() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storage_ClearScores_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClearScores'
type Storage_ClearScores_Call struct {
	*mock.Call
}

// ClearScores is a helper method to define mock.On call
func (_e *Storage_Expecter) ClearScores() *Storage_ClearScores_Call {
	return &Storage_ClearScores_Call{Call: _e.mock.On("ClearScores")}
}

func (_c *Storage_ClearScores_Call) Run(run func()) *Storage_ClearScores_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Storage_ClearScores_Call) Return(_a0 error) *Storage_ClearScores_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetLastScores provides a mock function with given fields:
func (_m *Storage) GetLastScores() ([]types.Result, error) {
	ret := _m.Called()

	var r0 []types.Result
	if rf, ok := ret.Get(0).(func() []types.Result); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storage_GetLastScores_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLastScores'
type Storage_GetLastScores_Call struct {
	*mock.Call
}

// GetLastScores is a helper method to define mock.On call
func (_e *Storage_Expecter) GetLastScores() *Storage_GetLastScores_Call {
	return &Storage_GetLastScores_Call{Call: _e.mock.On("GetLastScores")}
}

func (_c *Storage_GetLastScores_Call) Run(run func()) *Storage_GetLastScores_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Storage_GetLastScores_Call) Return(_a0 []types.Result, _a1 error) *Storage_GetLastScores_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// SetLastScore provides a mock function with given fields: _a0
func (_m *Storage) SetLastScore(_a0 types.Result) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Result) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storage_SetLastScore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetLastScore'
type Storage_SetLastScore_Call struct {
	*mock.Call
}

// SetLastScore is a helper method to define mock.On call
//   - _a0 types.Result
func (_e *Storage_Expecter) SetLastScore(_a0 interface{}) *Storage_SetLastScore_Call {
	return &Storage_SetLastScore_Call{Call: _e.mock.On("SetLastScore", _a0)}
}

func (_c *Storage_SetLastScore_Call) Run(run func(_a0 types.Result)) *Storage_SetLastScore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Result))
	})
	return _c
}

func (_c *Storage_SetLastScore_Call) Return(_a0 error) *Storage_SetLastScore_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	pkg "github.com/complynx/rpssl4bu/backend/pkg"
	mock "github.com/stretchr/testify/mock"

	types "github.com/complynx/rpssl4bu/backend/pkg/types"
)

// P2PGameFactory is an autogenerated mock type for the P2PGameFactory type
type P2PGameFactory struct {
	mock.Mock
}

type P2PGameFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *P2PGameFactory) EXPECT() *P2PGameFactory_Expecter {
	return &P2PGameFactory_Expecter{mock: &_m.Mock}
}

// CreateGame provides a mock function with given fields: ctx
func (_m *P2PGameFactory) CreateGame(ctx context.Context) (pkg.P2PGame, error) {
	ret := _m.Called(ctx)

	var r0 pkg.P2PGame
	if rf, ok := ret.Get(0).(func(context.Context) pkg.P2PGame); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pkg.P2PGame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// P2PGameFactory_CreateGame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateGame'
type P2PGameFactory_CreateGame_Call struct {
	*mock.Call
}

// CreateGame is a helper method to define mock.On call
//   - ctx context.Context
func (_e *P2PGameFactory_Expecter) CreateGame(ctx interface{}) *P2PGameFactory_CreateGame_Call {
	return &P2PGameFactory_CreateGame_Call{Call: _e.mock.On("CreateGame", ctx)}
}

func (_c *P2PGameFactory_CreateGame_Call) Run(run func(ctx context.Context)) *P2PGameFactory_CreateGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *P2PGameFactory_CreateGame_Call) Return(_a0 pkg.P2PGame, _a1 error) *P2PGameFactory_CreateGame_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetGame provides a mock function with given fields: id
func (_m *P2PGameFactory) GetGame(id types.GameID) (pkg.P2PGame, bool) {
	ret := _m.Called(id)

	var r0 pkg.P2PGame
	if rf, ok := ret.Get(0).(func(types.GameID) pkg.P2PGame); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pkg.P2PGame)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(types.GameID) bool); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// P2PGameFactory_GetGame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGame'
type P2PGameFactory_GetGame_Call struct {
	*mock.Call
}

// GetGame is a helper method to define mock.On call
//   - id types.GameID
func (_e *P2PGameFactory_Expecter) GetGame(id interface{}) *P2PGameFactory_GetGame_Call {
	return &P2PGameFactory_GetGame_Call{Call: _e.mock.On("GetGame", id)}
}

func (_c *P2PGameFactory_GetGame_Call) Run(run func(id types.GameID)) *P2PGameFactory_GetGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.GameID))
	})
	return _c
}

func (_c *P2PGameFactory_GetGame_Call) Return(_a0 pkg.P2PGame, _a1 bool) *P2PGameFactory_GetGame_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// StopGames provides a mock function with given fields: ctx
func (_m *P2PGameFactory) StopGames(ctx context.Context) {
	_m.Called(ctx)
}

// P2PGameFactory_StopGames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopGames'
type P2PGameFactory_StopGames_Call struct {
	*mock.Call
}

// StopGames is a helper method to define mock.On call
//   - ctx context.Context
func (_e *P2PGameFactory_Expecter) StopGames(ctx interface{}) *P2PGameFactory_StopGames_Call {
	return &P2PGameFactory_StopGames_Call{Call: _e.mock.On("StopGames", ctx)}
}

func (_c *P2PGameFactory_StopGames_Call) Run(run func(ctx context.Context)) *P2PGameFactory_StopGames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *P2PGameFactory_StopGames_Call) Return() *P2PGameFactory_StopGames_Call {
	_c.Call.Return()
	return _c
}

type mockConstructorTestingTNewP2PGameFactory interface {
	mock.TestingT
	Cleanup(func())
}

// NewP2PGameFactory creates a new instance of P2PGameFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewP2PGameFactory(t mockConstructorTestingTNewP2PGameFactory) *P2PGameFactory {
	mock := &P2PGameFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

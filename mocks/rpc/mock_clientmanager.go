// Code generated by mockery v2.43.0. DO NOT EDIT.

package rpc

import (
	mock "github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"

	rpc "github.com/evorts/kevlars/rpc"
)

// ClientManager is an autogenerated mock type for the ClientManager type
type ClientManager struct {
	mock.Mock
}

type ClientManager_Expecter struct {
	mock *mock.Mock
}

func (_m *ClientManager) EXPECT() *ClientManager_Expecter {
	return &ClientManager_Expecter{mock: &_m.Mock}
}

// Client provides a mock function with given fields:
func (_m *ClientManager) Client() grpc.ClientConnInterface {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Client")
	}

	var r0 grpc.ClientConnInterface
	if rf, ok := ret.Get(0).(func() grpc.ClientConnInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(grpc.ClientConnInterface)
		}
	}

	return r0
}

// ClientManager_Client_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Client'
type ClientManager_Client_Call struct {
	*mock.Call
}

// Client is a helper method to define mock.On call
func (_e *ClientManager_Expecter) Client() *ClientManager_Client_Call {
	return &ClientManager_Client_Call{Call: _e.mock.On("Client")}
}

func (_c *ClientManager_Client_Call) Run(run func()) *ClientManager_Client_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ClientManager_Client_Call) Return(_a0 grpc.ClientConnInterface) *ClientManager_Client_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientManager_Client_Call) RunAndReturn(run func() grpc.ClientConnInterface) *ClientManager_Client_Call {
	_c.Call.Return(run)
	return _c
}

// MustConnect provides a mock function with given fields:
func (_m *ClientManager) MustConnect() rpc.ClientManager {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MustConnect")
	}

	var r0 rpc.ClientManager
	if rf, ok := ret.Get(0).(func() rpc.ClientManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rpc.ClientManager)
		}
	}

	return r0
}

// ClientManager_MustConnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MustConnect'
type ClientManager_MustConnect_Call struct {
	*mock.Call
}

// MustConnect is a helper method to define mock.On call
func (_e *ClientManager_Expecter) MustConnect() *ClientManager_MustConnect_Call {
	return &ClientManager_MustConnect_Call{Call: _e.mock.On("MustConnect")}
}

func (_c *ClientManager_MustConnect_Call) Run(run func()) *ClientManager_MustConnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ClientManager_MustConnect_Call) Return(_a0 rpc.ClientManager) *ClientManager_MustConnect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientManager_MustConnect_Call) RunAndReturn(run func() rpc.ClientManager) *ClientManager_MustConnect_Call {
	_c.Call.Return(run)
	return _c
}

// Teardown provides a mock function with given fields:
func (_m *ClientManager) Teardown() {
	_m.Called()
}

// ClientManager_Teardown_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Teardown'
type ClientManager_Teardown_Call struct {
	*mock.Call
}

// Teardown is a helper method to define mock.On call
func (_e *ClientManager_Expecter) Teardown() *ClientManager_Teardown_Call {
	return &ClientManager_Teardown_Call{Call: _e.mock.On("Teardown")}
}

func (_c *ClientManager_Teardown_Call) Run(run func()) *ClientManager_Teardown_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ClientManager_Teardown_Call) Return() *ClientManager_Teardown_Call {
	_c.Call.Return()
	return _c
}

func (_c *ClientManager_Teardown_Call) RunAndReturn(run func()) *ClientManager_Teardown_Call {
	_c.Call.Return(run)
	return _c
}

// connect provides a mock function with given fields:
func (_m *ClientManager) connect() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for connect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClientManager_connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'connect'
type ClientManager_connect_Call struct {
	*mock.Call
}

// connect is a helper method to define mock.On call
func (_e *ClientManager_Expecter) connect() *ClientManager_connect_Call {
	return &ClientManager_connect_Call{Call: _e.mock.On("connect")}
}

func (_c *ClientManager_connect_Call) Run(run func()) *ClientManager_connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ClientManager_connect_Call) Return(_a0 error) *ClientManager_connect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientManager_connect_Call) RunAndReturn(run func() error) *ClientManager_connect_Call {
	_c.Call.Return(run)
	return _c
}

// NewClientManager creates a new instance of ClientManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClientManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClientManager {
	mock := &ClientManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

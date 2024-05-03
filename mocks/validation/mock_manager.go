// Code generated by mockery v2.42.3. DO NOT EDIT.

package validation

import (
	regexp "regexp"

	mock "github.com/stretchr/testify/mock"

	validation "github.com/evorts/kevlars/validation"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

type Manager_Expecter struct {
	mock *mock.Mock
}

func (_m *Manager) EXPECT() *Manager_Expecter {
	return &Manager_Expecter{mock: &_m.Mock}
}

// Engine provides a mock function with given fields:
func (_m *Manager) Engine() interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Engine")
	}

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Manager_Engine_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Engine'
type Manager_Engine_Call struct {
	*mock.Call
}

// Engine is a helper method to define mock.On call
func (_e *Manager_Expecter) Engine() *Manager_Engine_Call {
	return &Manager_Engine_Call{Call: _e.mock.On("Engine")}
}

func (_c *Manager_Engine_Call) Run(run func()) *Manager_Engine_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Manager_Engine_Call) Return(_a0 interface{}) *Manager_Engine_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Manager_Engine_Call) RunAndReturn(run func() interface{}) *Manager_Engine_Call {
	_c.Call.Return(run)
	return _c
}

// MustInit provides a mock function with given fields:
func (_m *Manager) MustInit() validation.Manager {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MustInit")
	}

	var r0 validation.Manager
	if rf, ok := ret.Get(0).(func() validation.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(validation.Manager)
		}
	}

	return r0
}

// Manager_MustInit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MustInit'
type Manager_MustInit_Call struct {
	*mock.Call
}

// MustInit is a helper method to define mock.On call
func (_e *Manager_Expecter) MustInit() *Manager_MustInit_Call {
	return &Manager_MustInit_Call{Call: _e.mock.On("MustInit")}
}

func (_c *Manager_MustInit_Call) Run(run func()) *Manager_MustInit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Manager_MustInit_Call) Return(_a0 validation.Manager) *Manager_MustInit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Manager_MustInit_Call) RunAndReturn(run func() validation.Manager) *Manager_MustInit_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterRegexValidator provides a mock function with given fields: tag, r
func (_m *Manager) RegisterRegexValidator(tag string, r *regexp.Regexp) validation.Manager {
	ret := _m.Called(tag, r)

	if len(ret) == 0 {
		panic("no return value specified for RegisterRegexValidator")
	}

	var r0 validation.Manager
	if rf, ok := ret.Get(0).(func(string, *regexp.Regexp) validation.Manager); ok {
		r0 = rf(tag, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(validation.Manager)
		}
	}

	return r0
}

// Manager_RegisterRegexValidator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterRegexValidator'
type Manager_RegisterRegexValidator_Call struct {
	*mock.Call
}

// RegisterRegexValidator is a helper method to define mock.On call
//   - tag string
//   - r *regexp.Regexp
func (_e *Manager_Expecter) RegisterRegexValidator(tag interface{}, r interface{}) *Manager_RegisterRegexValidator_Call {
	return &Manager_RegisterRegexValidator_Call{Call: _e.mock.On("RegisterRegexValidator", tag, r)}
}

func (_c *Manager_RegisterRegexValidator_Call) Run(run func(tag string, r *regexp.Regexp)) *Manager_RegisterRegexValidator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*regexp.Regexp))
	})
	return _c
}

func (_c *Manager_RegisterRegexValidator_Call) Return(_a0 validation.Manager) *Manager_RegisterRegexValidator_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Manager_RegisterRegexValidator_Call) RunAndReturn(run func(string, *regexp.Regexp) validation.Manager) *Manager_RegisterRegexValidator_Call {
	_c.Call.Return(run)
	return _c
}

// Validate provides a mock function with given fields: i
func (_m *Manager) Validate(i interface{}) error {
	ret := _m.Called(i)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Manager_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type Manager_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - i interface{}
func (_e *Manager_Expecter) Validate(i interface{}) *Manager_Validate_Call {
	return &Manager_Validate_Call{Call: _e.mock.On("Validate", i)}
}

func (_c *Manager_Validate_Call) Run(run func(i interface{})) *Manager_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Manager_Validate_Call) Return(_a0 error) *Manager_Validate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Manager_Validate_Call) RunAndReturn(run func(interface{}) error) *Manager_Validate_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateStruct provides a mock function with given fields: i
func (_m *Manager) ValidateStruct(i interface{}) error {
	ret := _m.Called(i)

	if len(ret) == 0 {
		panic("no return value specified for ValidateStruct")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Manager_ValidateStruct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateStruct'
type Manager_ValidateStruct_Call struct {
	*mock.Call
}

// ValidateStruct is a helper method to define mock.On call
//   - i interface{}
func (_e *Manager_Expecter) ValidateStruct(i interface{}) *Manager_ValidateStruct_Call {
	return &Manager_ValidateStruct_Call{Call: _e.mock.On("ValidateStruct", i)}
}

func (_c *Manager_ValidateStruct_Call) Run(run func(i interface{})) *Manager_ValidateStruct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Manager_ValidateStruct_Call) Return(_a0 error) *Manager_ValidateStruct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Manager_ValidateStruct_Call) RunAndReturn(run func(interface{}) error) *Manager_ValidateStruct_Call {
	_c.Call.Return(run)
	return _c
}

// NewManager creates a new instance of Manager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *Manager {
	mock := &Manager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

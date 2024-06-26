// Code generated by mockery v2.43.0. DO NOT EDIT.

package logger

import mock "github.com/stretchr/testify/mock"

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

// Debug provides a mock function with given fields: messages
func (_m *Manager) Debug(messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_Debug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debug'
type Manager_Debug_Call struct {
	*mock.Call
}

// Debug is a helper method to define mock.On call
//   - messages ...interface{}
func (_e *Manager_Expecter) Debug(messages ...interface{}) *Manager_Debug_Call {
	return &Manager_Debug_Call{Call: _e.mock.On("Debug",
		append([]interface{}{}, messages...)...)}
}

func (_c *Manager_Debug_Call) Run(run func(messages ...interface{})) *Manager_Debug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Manager_Debug_Call) Return() *Manager_Debug_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_Debug_Call) RunAndReturn(run func(...interface{})) *Manager_Debug_Call {
	_c.Call.Return(run)
	return _c
}

// DebugWhen provides a mock function with given fields: expr, f
func (_m *Manager) DebugWhen(expr bool, f func(func(...interface{}))) {
	_m.Called(expr, f)
}

// Manager_DebugWhen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DebugWhen'
type Manager_DebugWhen_Call struct {
	*mock.Call
}

// DebugWhen is a helper method to define mock.On call
//   - expr bool
//   - f func(func(...interface{}))
func (_e *Manager_Expecter) DebugWhen(expr interface{}, f interface{}) *Manager_DebugWhen_Call {
	return &Manager_DebugWhen_Call{Call: _e.mock.On("DebugWhen", expr, f)}
}

func (_c *Manager_DebugWhen_Call) Run(run func(expr bool, f func(func(...interface{})))) *Manager_DebugWhen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(func(func(...interface{}))))
	})
	return _c
}

func (_c *Manager_DebugWhen_Call) Return() *Manager_DebugWhen_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_DebugWhen_Call) RunAndReturn(run func(bool, func(func(...interface{})))) *Manager_DebugWhen_Call {
	_c.Call.Return(run)
	return _c
}

// DebugWithProps provides a mock function with given fields: props, messages
func (_m *Manager) DebugWithProps(props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_DebugWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DebugWithProps'
type Manager_DebugWithProps_Call struct {
	*mock.Call
}

// DebugWithProps is a helper method to define mock.On call
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) DebugWithProps(props interface{}, messages ...interface{}) *Manager_DebugWithProps_Call {
	return &Manager_DebugWithProps_Call{Call: _e.mock.On("DebugWithProps",
		append([]interface{}{props}, messages...)...)}
}

func (_c *Manager_DebugWithProps_Call) Run(run func(props map[string]interface{}, messages ...interface{})) *Manager_DebugWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_DebugWithProps_Call) Return() *Manager_DebugWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_DebugWithProps_Call) RunAndReturn(run func(map[string]interface{}, ...interface{})) *Manager_DebugWithProps_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields: messages
func (_m *Manager) Error(messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type Manager_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//   - messages ...interface{}
func (_e *Manager_Expecter) Error(messages ...interface{}) *Manager_Error_Call {
	return &Manager_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{}, messages...)...)}
}

func (_c *Manager_Error_Call) Run(run func(messages ...interface{})) *Manager_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Manager_Error_Call) Return() *Manager_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_Error_Call) RunAndReturn(run func(...interface{})) *Manager_Error_Call {
	_c.Call.Return(run)
	return _c
}

// ErrorWhen provides a mock function with given fields: expr, f
func (_m *Manager) ErrorWhen(expr bool, f func(func(...interface{}))) {
	_m.Called(expr, f)
}

// Manager_ErrorWhen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ErrorWhen'
type Manager_ErrorWhen_Call struct {
	*mock.Call
}

// ErrorWhen is a helper method to define mock.On call
//   - expr bool
//   - f func(func(...interface{}))
func (_e *Manager_Expecter) ErrorWhen(expr interface{}, f interface{}) *Manager_ErrorWhen_Call {
	return &Manager_ErrorWhen_Call{Call: _e.mock.On("ErrorWhen", expr, f)}
}

func (_c *Manager_ErrorWhen_Call) Run(run func(expr bool, f func(func(...interface{})))) *Manager_ErrorWhen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(func(func(...interface{}))))
	})
	return _c
}

func (_c *Manager_ErrorWhen_Call) Return() *Manager_ErrorWhen_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_ErrorWhen_Call) RunAndReturn(run func(bool, func(func(...interface{})))) *Manager_ErrorWhen_Call {
	_c.Call.Return(run)
	return _c
}

// ErrorWithProps provides a mock function with given fields: props, messages
func (_m *Manager) ErrorWithProps(props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_ErrorWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ErrorWithProps'
type Manager_ErrorWithProps_Call struct {
	*mock.Call
}

// ErrorWithProps is a helper method to define mock.On call
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) ErrorWithProps(props interface{}, messages ...interface{}) *Manager_ErrorWithProps_Call {
	return &Manager_ErrorWithProps_Call{Call: _e.mock.On("ErrorWithProps",
		append([]interface{}{props}, messages...)...)}
}

func (_c *Manager_ErrorWithProps_Call) Run(run func(props map[string]interface{}, messages ...interface{})) *Manager_ErrorWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_ErrorWithProps_Call) Return() *Manager_ErrorWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_ErrorWithProps_Call) RunAndReturn(run func(map[string]interface{}, ...interface{})) *Manager_ErrorWithProps_Call {
	_c.Call.Return(run)
	return _c
}

// ErrorWithPropsWhen provides a mock function with given fields: expr, props, f
func (_m *Manager) ErrorWithPropsWhen(expr bool, props map[string]interface{}, f func(func(...interface{}))) {
	_m.Called(expr, props, f)
}

// Manager_ErrorWithPropsWhen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ErrorWithPropsWhen'
type Manager_ErrorWithPropsWhen_Call struct {
	*mock.Call
}

// ErrorWithPropsWhen is a helper method to define mock.On call
//   - expr bool
//   - props map[string]interface{}
//   - f func(func(...interface{}))
func (_e *Manager_Expecter) ErrorWithPropsWhen(expr interface{}, props interface{}, f interface{}) *Manager_ErrorWithPropsWhen_Call {
	return &Manager_ErrorWithPropsWhen_Call{Call: _e.mock.On("ErrorWithPropsWhen", expr, props, f)}
}

func (_c *Manager_ErrorWithPropsWhen_Call) Run(run func(expr bool, props map[string]interface{}, f func(func(...interface{})))) *Manager_ErrorWithPropsWhen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(map[string]interface{}), args[2].(func(func(...interface{}))))
	})
	return _c
}

func (_c *Manager_ErrorWithPropsWhen_Call) Return() *Manager_ErrorWithPropsWhen_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_ErrorWithPropsWhen_Call) RunAndReturn(run func(bool, map[string]interface{}, func(func(...interface{})))) *Manager_ErrorWithPropsWhen_Call {
	_c.Call.Return(run)
	return _c
}

// Fatal provides a mock function with given fields: messages
func (_m *Manager) Fatal(messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_Fatal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fatal'
type Manager_Fatal_Call struct {
	*mock.Call
}

// Fatal is a helper method to define mock.On call
//   - messages ...interface{}
func (_e *Manager_Expecter) Fatal(messages ...interface{}) *Manager_Fatal_Call {
	return &Manager_Fatal_Call{Call: _e.mock.On("Fatal",
		append([]interface{}{}, messages...)...)}
}

func (_c *Manager_Fatal_Call) Run(run func(messages ...interface{})) *Manager_Fatal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Manager_Fatal_Call) Return() *Manager_Fatal_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_Fatal_Call) RunAndReturn(run func(...interface{})) *Manager_Fatal_Call {
	_c.Call.Return(run)
	return _c
}

// FatalWithProps provides a mock function with given fields: props, messages
func (_m *Manager) FatalWithProps(props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_FatalWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FatalWithProps'
type Manager_FatalWithProps_Call struct {
	*mock.Call
}

// FatalWithProps is a helper method to define mock.On call
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) FatalWithProps(props interface{}, messages ...interface{}) *Manager_FatalWithProps_Call {
	return &Manager_FatalWithProps_Call{Call: _e.mock.On("FatalWithProps",
		append([]interface{}{props}, messages...)...)}
}

func (_c *Manager_FatalWithProps_Call) Run(run func(props map[string]interface{}, messages ...interface{})) *Manager_FatalWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_FatalWithProps_Call) Return() *Manager_FatalWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_FatalWithProps_Call) RunAndReturn(run func(map[string]interface{}, ...interface{})) *Manager_FatalWithProps_Call {
	_c.Call.Return(run)
	return _c
}

// Info provides a mock function with given fields: messages
func (_m *Manager) Info(messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type Manager_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//   - messages ...interface{}
func (_e *Manager_Expecter) Info(messages ...interface{}) *Manager_Info_Call {
	return &Manager_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{}, messages...)...)}
}

func (_c *Manager_Info_Call) Run(run func(messages ...interface{})) *Manager_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Manager_Info_Call) Return() *Manager_Info_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_Info_Call) RunAndReturn(run func(...interface{})) *Manager_Info_Call {
	_c.Call.Return(run)
	return _c
}

// InfoWhen provides a mock function with given fields: expr, f
func (_m *Manager) InfoWhen(expr bool, f func(func(...interface{}))) {
	_m.Called(expr, f)
}

// Manager_InfoWhen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InfoWhen'
type Manager_InfoWhen_Call struct {
	*mock.Call
}

// InfoWhen is a helper method to define mock.On call
//   - expr bool
//   - f func(func(...interface{}))
func (_e *Manager_Expecter) InfoWhen(expr interface{}, f interface{}) *Manager_InfoWhen_Call {
	return &Manager_InfoWhen_Call{Call: _e.mock.On("InfoWhen", expr, f)}
}

func (_c *Manager_InfoWhen_Call) Run(run func(expr bool, f func(func(...interface{})))) *Manager_InfoWhen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(func(func(...interface{}))))
	})
	return _c
}

func (_c *Manager_InfoWhen_Call) Return() *Manager_InfoWhen_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_InfoWhen_Call) RunAndReturn(run func(bool, func(func(...interface{})))) *Manager_InfoWhen_Call {
	_c.Call.Return(run)
	return _c
}

// InfoWithProps provides a mock function with given fields: props, messages
func (_m *Manager) InfoWithProps(props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_InfoWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InfoWithProps'
type Manager_InfoWithProps_Call struct {
	*mock.Call
}

// InfoWithProps is a helper method to define mock.On call
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) InfoWithProps(props interface{}, messages ...interface{}) *Manager_InfoWithProps_Call {
	return &Manager_InfoWithProps_Call{Call: _e.mock.On("InfoWithProps",
		append([]interface{}{props}, messages...)...)}
}

func (_c *Manager_InfoWithProps_Call) Run(run func(props map[string]interface{}, messages ...interface{})) *Manager_InfoWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_InfoWithProps_Call) Return() *Manager_InfoWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_InfoWithProps_Call) RunAndReturn(run func(map[string]interface{}, ...interface{})) *Manager_InfoWithProps_Call {
	_c.Call.Return(run)
	return _c
}

// InfoWithPropsWhen provides a mock function with given fields: expr, props, messages
func (_m *Manager) InfoWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, expr, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_InfoWithPropsWhen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InfoWithPropsWhen'
type Manager_InfoWithPropsWhen_Call struct {
	*mock.Call
}

// InfoWithPropsWhen is a helper method to define mock.On call
//   - expr bool
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) InfoWithPropsWhen(expr interface{}, props interface{}, messages ...interface{}) *Manager_InfoWithPropsWhen_Call {
	return &Manager_InfoWithPropsWhen_Call{Call: _e.mock.On("InfoWithPropsWhen",
		append([]interface{}{expr, props}, messages...)...)}
}

func (_c *Manager_InfoWithPropsWhen_Call) Run(run func(expr bool, props map[string]interface{}, messages ...interface{})) *Manager_InfoWithPropsWhen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(bool), args[1].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_InfoWithPropsWhen_Call) Return() *Manager_InfoWithPropsWhen_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_InfoWithPropsWhen_Call) RunAndReturn(run func(bool, map[string]interface{}, ...interface{})) *Manager_InfoWithPropsWhen_Call {
	_c.Call.Return(run)
	return _c
}

// Panic provides a mock function with given fields: messages
func (_m *Manager) Panic(messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_Panic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Panic'
type Manager_Panic_Call struct {
	*mock.Call
}

// Panic is a helper method to define mock.On call
//   - messages ...interface{}
func (_e *Manager_Expecter) Panic(messages ...interface{}) *Manager_Panic_Call {
	return &Manager_Panic_Call{Call: _e.mock.On("Panic",
		append([]interface{}{}, messages...)...)}
}

func (_c *Manager_Panic_Call) Run(run func(messages ...interface{})) *Manager_Panic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Manager_Panic_Call) Return() *Manager_Panic_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_Panic_Call) RunAndReturn(run func(...interface{})) *Manager_Panic_Call {
	_c.Call.Return(run)
	return _c
}

// PanicWithProps provides a mock function with given fields: props, messages
func (_m *Manager) PanicWithProps(props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_PanicWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PanicWithProps'
type Manager_PanicWithProps_Call struct {
	*mock.Call
}

// PanicWithProps is a helper method to define mock.On call
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) PanicWithProps(props interface{}, messages ...interface{}) *Manager_PanicWithProps_Call {
	return &Manager_PanicWithProps_Call{Call: _e.mock.On("PanicWithProps",
		append([]interface{}{props}, messages...)...)}
}

func (_c *Manager_PanicWithProps_Call) Run(run func(props map[string]interface{}, messages ...interface{})) *Manager_PanicWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_PanicWithProps_Call) Return() *Manager_PanicWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_PanicWithProps_Call) RunAndReturn(run func(map[string]interface{}, ...interface{})) *Manager_PanicWithProps_Call {
	_c.Call.Return(run)
	return _c
}

// Trace provides a mock function with given fields: messages
func (_m *Manager) Trace(messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_Trace_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Trace'
type Manager_Trace_Call struct {
	*mock.Call
}

// Trace is a helper method to define mock.On call
//   - messages ...interface{}
func (_e *Manager_Expecter) Trace(messages ...interface{}) *Manager_Trace_Call {
	return &Manager_Trace_Call{Call: _e.mock.On("Trace",
		append([]interface{}{}, messages...)...)}
}

func (_c *Manager_Trace_Call) Run(run func(messages ...interface{})) *Manager_Trace_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Manager_Trace_Call) Return() *Manager_Trace_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_Trace_Call) RunAndReturn(run func(...interface{})) *Manager_Trace_Call {
	_c.Call.Return(run)
	return _c
}

// TraceWithProps provides a mock function with given fields: props, messages
func (_m *Manager) TraceWithProps(props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_TraceWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TraceWithProps'
type Manager_TraceWithProps_Call struct {
	*mock.Call
}

// TraceWithProps is a helper method to define mock.On call
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) TraceWithProps(props interface{}, messages ...interface{}) *Manager_TraceWithProps_Call {
	return &Manager_TraceWithProps_Call{Call: _e.mock.On("TraceWithProps",
		append([]interface{}{props}, messages...)...)}
}

func (_c *Manager_TraceWithProps_Call) Run(run func(props map[string]interface{}, messages ...interface{})) *Manager_TraceWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_TraceWithProps_Call) Return() *Manager_TraceWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_TraceWithProps_Call) RunAndReturn(run func(map[string]interface{}, ...interface{})) *Manager_TraceWithProps_Call {
	_c.Call.Return(run)
	return _c
}

// Warn provides a mock function with given fields: messages
func (_m *Manager) Warn(messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_Warn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warn'
type Manager_Warn_Call struct {
	*mock.Call
}

// Warn is a helper method to define mock.On call
//   - messages ...interface{}
func (_e *Manager_Expecter) Warn(messages ...interface{}) *Manager_Warn_Call {
	return &Manager_Warn_Call{Call: _e.mock.On("Warn",
		append([]interface{}{}, messages...)...)}
}

func (_c *Manager_Warn_Call) Run(run func(messages ...interface{})) *Manager_Warn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Manager_Warn_Call) Return() *Manager_Warn_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_Warn_Call) RunAndReturn(run func(...interface{})) *Manager_Warn_Call {
	_c.Call.Return(run)
	return _c
}

// WarnWhen provides a mock function with given fields: expr, f
func (_m *Manager) WarnWhen(expr bool, f func(func(...interface{}))) {
	_m.Called(expr, f)
}

// Manager_WarnWhen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WarnWhen'
type Manager_WarnWhen_Call struct {
	*mock.Call
}

// WarnWhen is a helper method to define mock.On call
//   - expr bool
//   - f func(func(...interface{}))
func (_e *Manager_Expecter) WarnWhen(expr interface{}, f interface{}) *Manager_WarnWhen_Call {
	return &Manager_WarnWhen_Call{Call: _e.mock.On("WarnWhen", expr, f)}
}

func (_c *Manager_WarnWhen_Call) Run(run func(expr bool, f func(func(...interface{})))) *Manager_WarnWhen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(func(func(...interface{}))))
	})
	return _c
}

func (_c *Manager_WarnWhen_Call) Return() *Manager_WarnWhen_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_WarnWhen_Call) RunAndReturn(run func(bool, func(func(...interface{})))) *Manager_WarnWhen_Call {
	_c.Call.Return(run)
	return _c
}

// WarnWithProps provides a mock function with given fields: props, messages
func (_m *Manager) WarnWithProps(props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_WarnWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WarnWithProps'
type Manager_WarnWithProps_Call struct {
	*mock.Call
}

// WarnWithProps is a helper method to define mock.On call
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) WarnWithProps(props interface{}, messages ...interface{}) *Manager_WarnWithProps_Call {
	return &Manager_WarnWithProps_Call{Call: _e.mock.On("WarnWithProps",
		append([]interface{}{props}, messages...)...)}
}

func (_c *Manager_WarnWithProps_Call) Run(run func(props map[string]interface{}, messages ...interface{})) *Manager_WarnWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_WarnWithProps_Call) Return() *Manager_WarnWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_WarnWithProps_Call) RunAndReturn(run func(map[string]interface{}, ...interface{})) *Manager_WarnWithProps_Call {
	_c.Call.Return(run)
	return _c
}

// WarnWithPropsWhen provides a mock function with given fields: expr, props, messages
func (_m *Manager) WarnWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, expr, props)
	_ca = append(_ca, messages...)
	_m.Called(_ca...)
}

// Manager_WarnWithPropsWhen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WarnWithPropsWhen'
type Manager_WarnWithPropsWhen_Call struct {
	*mock.Call
}

// WarnWithPropsWhen is a helper method to define mock.On call
//   - expr bool
//   - props map[string]interface{}
//   - messages ...interface{}
func (_e *Manager_Expecter) WarnWithPropsWhen(expr interface{}, props interface{}, messages ...interface{}) *Manager_WarnWithPropsWhen_Call {
	return &Manager_WarnWithPropsWhen_Call{Call: _e.mock.On("WarnWithPropsWhen",
		append([]interface{}{expr, props}, messages...)...)}
}

func (_c *Manager_WarnWithPropsWhen_Call) Run(run func(expr bool, props map[string]interface{}, messages ...interface{})) *Manager_WarnWithPropsWhen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(bool), args[1].(map[string]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *Manager_WarnWithPropsWhen_Call) Return() *Manager_WarnWithPropsWhen_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_WarnWithPropsWhen_Call) RunAndReturn(run func(bool, map[string]interface{}, ...interface{})) *Manager_WarnWithPropsWhen_Call {
	_c.Call.Return(run)
	return _c
}

// WhenError provides a mock function with given fields: err
func (_m *Manager) WhenError(err error) {
	_m.Called(err)
}

// Manager_WhenError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WhenError'
type Manager_WhenError_Call struct {
	*mock.Call
}

// WhenError is a helper method to define mock.On call
//   - err error
func (_e *Manager_Expecter) WhenError(err interface{}) *Manager_WhenError_Call {
	return &Manager_WhenError_Call{Call: _e.mock.On("WhenError", err)}
}

func (_c *Manager_WhenError_Call) Run(run func(err error)) *Manager_WhenError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error))
	})
	return _c
}

func (_c *Manager_WhenError_Call) Return() *Manager_WhenError_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_WhenError_Call) RunAndReturn(run func(error)) *Manager_WhenError_Call {
	_c.Call.Return(run)
	return _c
}

// WhenErrorWithProps provides a mock function with given fields: err, props
func (_m *Manager) WhenErrorWithProps(err error, props map[string]interface{}) {
	_m.Called(err, props)
}

// Manager_WhenErrorWithProps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WhenErrorWithProps'
type Manager_WhenErrorWithProps_Call struct {
	*mock.Call
}

// WhenErrorWithProps is a helper method to define mock.On call
//   - err error
//   - props map[string]interface{}
func (_e *Manager_Expecter) WhenErrorWithProps(err interface{}, props interface{}) *Manager_WhenErrorWithProps_Call {
	return &Manager_WhenErrorWithProps_Call{Call: _e.mock.On("WhenErrorWithProps", err, props)}
}

func (_c *Manager_WhenErrorWithProps_Call) Run(run func(err error, props map[string]interface{})) *Manager_WhenErrorWithProps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(error), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *Manager_WhenErrorWithProps_Call) Return() *Manager_WhenErrorWithProps_Call {
	_c.Call.Return()
	return _c
}

func (_c *Manager_WhenErrorWithProps_Call) RunAndReturn(run func(error, map[string]interface{})) *Manager_WhenErrorWithProps_Call {
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

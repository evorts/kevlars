/**
 * @Author: steven
 * @Description:
 * @File: logger_noop
 * @Date: 18/12/23 00.51
 */

package logger

type noop struct{}

func (m *noop) Trace(messages ...interface{}) {
	// do nothing
}

func (m *noop) TraceWithProps(props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) Debug(messages ...interface{}) {
	// do nothing
}

func (m *noop) DebugWhen(expr bool, f func(debug func(messages ...interface{}))) {
	// do nothing
}

func (m *noop) DebugWithProps(props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) Info(messages ...interface{}) {
	// do nothing
}

func (m *noop) InfoWhen(expr bool, f func(info func(messages ...interface{}))) {
	// do nothing
}

func (m *noop) InfoWithProps(props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) InfoWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) Warn(messages ...interface{}) {
	// do nothing
}

func (m *noop) WarnWhen(expr bool, f func(warn func(messages ...interface{}))) {
	// do nothing
}

func (m *noop) WarnWithProps(props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) WarnWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) WhenError(err error) {
	// do nothing
}

func (m *noop) WhenErrorWithProps(err error, props map[string]interface{}) {
	// do nothing
}

func (m *noop) Error(messages ...interface{}) {
	// do nothing
}

func (m *noop) ErrorWhen(expr bool, f func(e func(messages ...interface{}))) {
	// do nothing
}

func (m *noop) ErrorWithProps(props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) ErrorWithPropsWhen(expr bool, props map[string]interface{}, f func(messages func(args ...interface{}))) {
	// do nothing
}

func (m *noop) Fatal(messages ...interface{}) {
	// do nothing
}

func (m *noop) FatalWithProps(props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func (m *noop) Panic(messages ...interface{}) {
	// do nothing
}

func (m *noop) PanicWithProps(props map[string]interface{}, messages ...interface{}) {
	// do nothing
}

func NewNoop() Manager {
	return &noop{}
}

/**
 * @Author: steven
 * @Description:
 * @File: logger
 * @Date: 18/12/23 00.50
 */

package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

// LogLevel 1 = DEBUG, 2 = INFO, 3 = WARN, 4 = ERROR, 5 = OFF, 6 = PANIC, 7 = FATAL
type LogLevel int

const (
	LogLevelDebug = 1
	LogLevelInfo  = 2
	LogLevelWarn  = 3
	LogLevelError = 4
	LogLevelOFF   = 5
	LogLevelPanic = 6
	LogLevelFatal = 7
)

func ParseLevel(value string) LogLevel {
	l, err := logrus.ParseLevel(value)
	if err != nil {
		return LogLevelError
	}
	switch l {
	case logrus.DebugLevel:
		return LogLevelDebug
	case logrus.InfoLevel:
		return LogLevelInfo
	case logrus.WarnLevel:
		return LogLevelWarn
	case logrus.TraceLevel:
		return LogLevelDebug
	case logrus.FatalLevel:
		return LogLevelFatal
	case logrus.ErrorLevel:
		fallthrough
	default:
		return LogLevelError
	}
}

func (l LogLevel) LogRushString() string {
	switch l {
	case LogLevelDebug:
		return logrus.DebugLevel.String()
	case LogLevelInfo:
		return logrus.InfoLevel.String()
	case LogLevelWarn:
		return logrus.WarnLevel.String()
	case LogLevelOFF:
		fallthrough
	case LogLevelPanic:
		return logrus.PanicLevel.String()
	case LogLevelFatal:
		return logrus.FatalLevel.String()
	case LogLevelError:
		fallthrough
	default:
		return logrus.ErrorLevel.String()
	}
}

type Manager interface {
	Trace(messages ...interface{})
	TraceWithProps(props map[string]interface{}, messages ...interface{})

	Debug(messages ...interface{})
	DebugWhen(expr bool, f func(debug func(messages ...interface{})))
	DebugWithProps(props map[string]interface{}, messages ...interface{})

	Info(messages ...interface{})
	InfoWhen(expr bool, f func(info func(messages ...interface{})))
	InfoWithProps(props map[string]interface{}, messages ...interface{})
	InfoWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{})

	Warn(messages ...interface{})
	WarnWhen(expr bool, f func(warn func(messages ...interface{})))
	WarnWithProps(props map[string]interface{}, messages ...interface{})
	WarnWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{})

	WhenError(err error)
	WhenErrorWithProps(err error, props map[string]interface{})

	Error(messages ...interface{})
	ErrorWhen(expr bool, f func(e func(messages ...interface{})))
	ErrorWithProps(props map[string]interface{}, messages ...interface{})
	ErrorWithPropsWhen(expr bool, props map[string]interface{}, f func(messages func(args ...interface{})))

	Fatal(messages ...interface{})
	FatalWithProps(props map[string]interface{}, messages ...interface{})

	Panic(messages ...interface{})
	PanicWithProps(props map[string]interface{}, messages ...interface{})
}

type manager struct {
	l    *logrus.Logger
	f    logrus.Formatter
	name string
}

func (m *manager) Trace(messages ...interface{}) {
	m.l.Traceln(messages...)
}

func (m *manager) TraceWithProps(props map[string]interface{}, messages ...interface{}) {
	m.l.WithFields(m.buildProps(props)).Traceln(messages...)
}

func (m *manager) Debug(messages ...interface{}) {
	m.l.Debugln(messages...)
}

func (m *manager) DebugWhen(expr bool, f func(debug func(messages ...interface{}))) {
	if expr {
		f(m.Debug)
	}
}

func (m *manager) DebugWithProps(props map[string]interface{}, messages ...interface{}) {
	m.l.WithFields(m.buildProps(props)).Debugln(messages...)
}

func (m *manager) Info(messages ...interface{}) {
	m.l.Infoln(messages...)
}

func (m *manager) InfoWhen(expr bool, f func(messages func(...interface{}))) {
	if expr {
		f(m.Info)
	}
}

func (m *manager) InfoWithProps(props map[string]interface{}, messages ...interface{}) {
	m.l.WithFields(m.buildProps(props)).Infoln(messages...)
}

func (m *manager) InfoWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{}) {
	if !expr {
		return
	}
	m.InfoWithProps(m.buildProps(props), messages...)
}

func (m *manager) Warn(messages ...interface{}) {
	m.l.Warnln(messages...)
}

func (m *manager) WarnWhen(expr bool, f func(messages func(...interface{}))) {
	if expr {
		f(m.Warn)
	}
}

func (m *manager) WarnWithProps(props map[string]interface{}, messages ...interface{}) {
	m.l.WithFields(m.buildProps(props)).Warnln(messages...)
}

func (m *manager) WarnWithPropsWhen(expr bool, props map[string]interface{}, messages ...interface{}) {
	if !expr {
		return
	}
	m.WarnWithProps(m.buildProps(props), messages...)
}

func (m *manager) ifErrorCallback(err error, callback func()) {
	if err != nil {
		callback()
	}
}
func (m *manager) WhenError(err error) {
	m.ifErrorCallback(err, func() {
		m.l.Errorln(err.Error())
	})
}

func (m *manager) WhenErrorWithProps(err error, props map[string]interface{}) {
	m.ifErrorCallback(err, func() {
		m.l.WithFields(m.buildProps(props)).Errorln(err.Error())
	})
}

func (m *manager) ErrorWhen(expr bool, f func(messages func(...interface{}))) {
	if !expr {
		return
	}
	f(m.l.Errorln)
}

func (m *manager) ErrorWithPropsWhen(expr bool, props map[string]interface{}, f func(messages func(...interface{}))) {
	if !expr {
		return
	}
	f(m.l.WithFields(m.buildProps(props)).Errorln)
}

func (m *manager) Error(messages ...interface{}) {
	m.l.Errorln(messages...)
}

func (m *manager) ErrorWithProps(props map[string]interface{}, messages ...interface{}) {
	m.l.WithFields(m.buildProps(props)).Errorln(messages...)
}

func (m *manager) Fatal(messages ...interface{}) {
	m.l.Fatal(messages...)
}

func (m *manager) FatalWithProps(props map[string]interface{}, messages ...interface{}) {
	m.l.WithFields(m.buildProps(props)).Fatalln(messages...)
}

func (m *manager) Panic(messages ...interface{}) {
	m.l.Panicln(messages...)
}

func (m *manager) PanicWithProps(props map[string]interface{}, messages ...interface{}) {
	m.l.WithFields(m.buildProps(props)).Panicln(messages...)
}

func (m *manager) buildProps(props map[string]interface{}) map[string]interface{} {
	props["app_name"] = m.name
	return props
}

func NewLogger(level LogLevel, output io.Writer, opts ...Option) Manager {
	l := logrus.New()
	logrusLogLevel, err := logrus.ParseLevel(level.LogRushString())
	if err != nil {
		logrusLogLevel = logrus.ErrorLevel
	}
	l.SetLevel(logrusLogLevel)
	l.SetOutput(output)
	l.SetFormatter(&logrus.JSONFormatter{})
	m := &manager{l: l}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}

/**
 * @Author: steven
 * @Description:
 * @File: logger_mock
 * @Date: 18/12/23 00.51
 */

package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

func NewMockLogger() Manager {
	l := logrus.New()
	l.SetLevel(logrus.ErrorLevel)
	l.SetOutput(io.Discard)
	l.SetFormatter(&logrus.JSONFormatter{})
	m := &manager{l: l}
	return m
}

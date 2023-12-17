/**
 * @Author: steven
 * @Description:
 * @File: tz_formatter
 * @Version: 1.0.0
 * @Date: 22/08/23 12.22
 */

package logger

import (
	"github.com/evorts/kevlars/ts"
	"github.com/sirupsen/logrus"
	"time"
)

type TZFormatter struct {
	tz  ts.TimeZone
	loc *time.Location
	logrus.Formatter
}

func (t TZFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if t.loc != nil {
		entry.Time = entry.Time.In(t.loc)
	}
	return t.Formatter.Format(entry)
}

func newTZFormatter(tz ts.TimeZone) (logrus.Formatter, error) {
	f := &TZFormatter{
		tz:        tz,
		Formatter: &logrus.JSONFormatter{},
	}
	var err error
	f.loc, err = time.LoadLocation(tz.String())
	return f, err
}

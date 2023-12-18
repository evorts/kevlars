package validation

import (
	"time"
)

const dateOnlyFormat = "2006-01-02"

// DateOnly expected format yyyy-mm-dd
type DateOnly string

func (d DateOnly) ToTime() *time.Time {
	t, err := time.Parse(dateOnlyFormat, d.String())
	if err != nil {
		return nil
	}
	return &t
}

func (d DateOnly) String() string {
	return string(d)
}

func ToDateOnly(t *time.Time) DateOnly {
	return DateOnly(t.Format(dateOnlyFormat))
}

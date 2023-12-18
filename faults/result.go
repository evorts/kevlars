/**
 * @Author: steven
 * @Description:
 * @File: result
 * @Date: 20/12/23 16.40
 */

package faults

import (
	"errors"
	"strings"
)

type ResultError struct {
	Errors []string
}

func (e *ResultError) Error() string {
	return strings.Join(e.Errors, "; ")
}

func (e *ResultError) Is(err error) bool {
	var resultError *ResultError
	ok := errors.As(err, &resultError)
	return ok
}

func (e *ResultError) Add(msg string) {
	e.Errors = append(e.Errors, msg)
}

func (e *ResultError) HasError() bool {
	return len(e.Errors) > 0
}

func (e *ResultError) ToError() error {
	if e.HasError() {
		return e
	}
	return nil
}

func NewResultError() *ResultError {
	return &ResultError{Errors: make([]string, 0)}
}

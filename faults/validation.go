/**
 * @Author: steven
 * @Description:
 * @File: validation
 * @Date: 15/11/23 15.59
 */

package faults

import "strings"

type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

func (e *ValidationError) Add(msg string) {
	e.Errors = append(e.Errors, msg)
}

func (e *ValidationError) HasError() bool {
	return len(e.Errors) > 0
}

func (e *ValidationError) ToError() error {
	if e.HasError() {
		return e
	}
	return nil
}

func NewValidationError() *ValidationError {
	return &ValidationError{Errors: make([]string, 0)}
}

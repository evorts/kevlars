/**
 * @Author: steven
 * @Description:
 * @File: validator
 * @Date: 29/09/23 10.50
 */

package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"sync"
	"time"
)

type Manager interface {
	Validate(i interface{}) error
	ValidateStruct(i interface{}) error
	RegisterRegexValidator(tag string, r *regexp.Regexp) Manager
	MustInit() Manager
	Engine() any
}

type manager struct {
	once      sync.Once
	validator *validator.Validate
}

func (v *manager) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func (v *manager) ValidateStruct(i interface{}) error {
	return v.validator.Struct(i)
}

func (v *manager) RegisterRegexValidator(tag string, r *regexp.Regexp) Manager {
	_ = v.validator.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		return r.MatchString(fl.Field().String())
	})
	return v
}

func (v *manager) MustInit() Manager {
	_ = v.validator.RegisterValidation("dateonly", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if !dateFormatRegex.MatchString(value) {
			return false
		}
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return false
		}
		return true
	})
	_ = v.validator.RegisterValidation("alphanumwithspace", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return alphaNumericWithSpace.MatchString(value)
	})
	_ = v.validator.RegisterValidation("alphanumwithdashandspace", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return alphaNumericWithDashAndSpace.MatchString(value)
	})
	_ = v.validator.RegisterValidation("datetime_rfc3339", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		_, err := time.Parse(time.RFC3339, value)
		return err == nil
	})
	return v
}

func (v *manager) lazyInit() {
	v.once.Do(func() {
		v.validator = validator.New()
		v.validator.SetTagName("validate")
	})
}

func (v *manager) Engine() any {
	return v.validator
}

func NewValidator() Manager {
	m := &manager{}
	m.lazyInit()
	return m
}

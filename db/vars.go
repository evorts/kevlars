/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 17/12/23 22.09
 */

package db

import "github.com/evorts/kevlars/utils"

const (
	DefaultOffset = 0
	DefaultPage   = 1
	DefaultLimit  = 10
)

type Operator string

const (
	OpEq    Operator = "="
	OpNotEq Operator = "<>"
	OpGt    Operator = ">"
	OpLt    Operator = "<"
	OpGte   Operator = ">="
	OpLte   Operator = "<="
	OpLike  Operator = "LIKE"
)

type Separator string

const (
	SeparatorAND Separator = "AND"
	SeparatorOR  Separator = "OR"
)

func (s Separator) String() string {
	return string(s)
}

func (s Separator) StringWithSpace() string {
	return " " + string(s) + " "
}

type Sort string

const (
	SortAsc  Sort = "asc"
	SortDesc Sort = "desc"
	SortNone Sort = "none"
)

func (t Sort) String() string {
	return string(t)
}

func (t Sort) Valid() bool {
	return utils.InArray([]string{SortAsc.String(), SortDesc.String(), SortNone.String()}, t.String())
}

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) ErrorCode() int {
	return e.Code
}

func NewError(code int, message string) Error {
	return Error{Code: code, Message: message}
}

var (
	ErrorRecordNotFound      error = NewError(404, "record not found")
	ErrorRecordAlreadyExists error = NewError(409, "record already exists")
)

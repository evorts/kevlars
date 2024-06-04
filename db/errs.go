/**
 * @Author: steven
 * @Description:
 * @File: errs
 * @Date: 02/06/24 21.44
 */

package db

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
	ErrorRecordNotFound      error = NewError(4040, "record not found")
	ErrorRecordAlreadyExists error = NewError(4090, "record already exists")
	ErrorEmptyArguments      error = NewError(4000, "arguments are empty")
	ErrorDriverNotSupported  error = NewError(4001, "driver not supported yet")
)

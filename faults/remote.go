/**
 * @Author: steven
 * @Description:
 * @File: custom
 * @Version: 1.0.0
 * @Date: 18/06/23 18.46
 */

package faults

type IRemoteError interface {
	IError
}

type remoteError struct {
	*commonError
}

func (e *remoteError) Error() string {
	return e.message
}

func (e *remoteError) Code() string {
	return e.code
}

func (e *remoteError) Props() map[string]string {
	return e.details
}

func (e *remoteError) UnWrap() error {
	return e.underlyingErr
}

func NewRemoteError(message string, opts ...IOption[remoteError]) IError {
	re := &remoteError{
		commonError: &commonError{
			message:       message,
			underlyingErr: nil,
		},
	}
	for _, opt := range opts {
		opt.apply(re)
	}
	return re
}

func NewRemoteErrorWrap(err error, message string, opts ...IOption[remoteError]) IError {
	re := &remoteError{
		commonError: &commonError{
			message:       message,
			underlyingErr: err,
		},
	}
	for _, opt := range opts {
		opt.apply(re)
	}
	return re
}

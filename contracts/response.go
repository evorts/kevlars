/**
 * @Author: steven
 * @Description:
 * @File: response
 * @Date: 22/12/23 10.52
 */

package contracts

import (
	"net/http"
)

type ResponseSuccess[T any] struct {
	HttpCode int    `json:"http_code"`
	Message  string `json:"message"`
	Content  T      `json:"content"`
}

type ErrorDetail struct {
	Code   string            `json:"code"`
	Errors map[string]string `json:"errors"`
}

type ResponseFail struct {
	HttpCode int         `json:"http_code"`
	Message  string      `json:"message"`
	Details  ErrorDetail `json:"details"`
}

func NewResponseSuccess[T any](httpCode int, message string, content T) (int, ResponseSuccess[T]) {
	return httpCode, ResponseSuccess[T]{
		HttpCode: httpCode,
		Message:  message,
		Content:  content,
	}
}

func NewResponseOK[T any](message string, content T) (int, ResponseSuccess[T]) {
	return http.StatusOK, ResponseSuccess[T]{
		HttpCode: http.StatusOK,
		Message:  message,
		Content:  content,
	}
}

func NewResponseFail(httpCode int, message string, details ErrorDetail) (int, ResponseFail) {
	return httpCode, ResponseFail{
		HttpCode: httpCode,
		Message:  message,
		Details:  details,
	}
}
func NewResponseBadRequest(message string) (int, ResponseFail) {
	return http.StatusBadRequest, ResponseFail{
		HttpCode: http.StatusBadRequest,
		Message:  message,
		Details:  ErrorDetail{},
	}
}
func NewResponseInternalServerError(message string) (int, ResponseFail) {
	return http.StatusInternalServerError, ResponseFail{
		HttpCode: http.StatusInternalServerError,
		Message:  message,
		Details:  ErrorDetail{},
	}
}

func NewResponseBadRequestWithDetail(message, errorCode string, errors map[string]string) (int, ResponseFail) {
	if errors == nil {
		errors = make(map[string]string)
	}
	return http.StatusBadRequest, ResponseFail{
		HttpCode: http.StatusBadRequest,
		Message:  message,
		Details: ErrorDetail{
			Code:   errorCode,
			Errors: errors,
		},
	}
}

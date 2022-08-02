package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/vonmutinda/organono/app/logger"
)

type ErrorCode string

var (
	ErrorCodeInvalidArgument    ErrorCode = "invalid_argument"
	ErrorCodeInvalidCredentials ErrorCode = "invalid_credentials"
	ErrorCodeInvalidForm        ErrorCode = "invalid_form"
	ErrorCodeInvalidPhone       ErrorCode = "invalid_phone"
	ErrorCodeInvalidUserStatus  ErrorCode = "invalid_user_status"
	ErrorCodeNotFound           ErrorCode = "not_found"
	ErrorCodeResourceExists     ErrorCode = "resource_exists"
	ErrorCodeRequestFailed      ErrorCode = "request_failed"
	ErrorCodeRoleForbidden      ErrorCode = "role_forbidden"
	ErrorCodeSessionExpired     ErrorCode = "session_expired"

	errorCodeMessageMap = map[ErrorCode]string{
		ErrorCodeInvalidArgument:    "You have provided an invalid argument",
		ErrorCodeInvalidCredentials: "You have provided invalid credentials",
		ErrorCodeInvalidForm:        "You have submitted an invalid form",
		ErrorCodeInvalidPhone:       "You have provided an invalid phone number",
		ErrorCodeInvalidUserStatus:  "Your account is not active",
		ErrorCodeNotFound:           "The requested resource was not found",
		ErrorCodeResourceExists:     "Another resource with similar attributes already exists",
		ErrorCodeRequestFailed:      "Request failed to complete. Please try again",
		ErrorCodeRoleForbidden:      "You are not allowed to perform this request",
		ErrorCodeSessionExpired:     "Your session has expired. Login again to proceed.",
	}

	httpStatusErrorCodeMap = map[ErrorCode]int{
		ErrorCodeInvalidCredentials: http.StatusUnauthorized,
		ErrorCodeInvalidUserStatus:  http.StatusNotAcceptable,
		ErrorCodeRoleForbidden:      http.StatusForbidden,
		ErrorCodeSessionExpired:     http.StatusUnauthorized,
	}
)

type Error struct {
	err            error
	errorCode      ErrorCode
	httpStatusCode int
	logMessages    []string
	ctx            context.Context
	notify         bool
}

func NewError(err error, format string, args ...interface{}) *Error {

	genericError, ok := err.(*Error)
	if !ok {
		return NewErrorWithCode(
			err,
			ErrorCodeRequestFailed,
			format,
			args...,
		)
	}

	genericError.addLogMessage(err, format, args...)

	return genericError
}

func NewErrorWithCode(
	err error,
	errorCode ErrorCode,
	format string,
	args ...interface{},
) *Error {

	if errorCode == "" {
		errorCode = ErrorCodeRequestFailed
	}

	genericError, ok := err.(*Error)
	if !ok {
		genericError = &Error{
			err: err,
		}
	}

	genericError.errorCode = errorCode

	genericError.addLogMessage(err, format, args...)

	return genericError
}

func (e *Error) Error() string {
	return strings.Join(e.logMessages, "; ")
}

func (e *Error) Err() error {
	return e.err
}

func (e *Error) GetErrorCode() ErrorCode {
	return e.errorCode
}

func (e *Error) JsonResponse() map[string]string {
	return map[string]string{
		"error_code":    e.errorCode.String(),
		"error_message": errorCodeMessageMap[e.errorCode],
	}
}

func (e *Error) WithContext(ctx context.Context) {
	e.ctx = ctx
}

func (e *Error) addLogMessage(
	err error,
	format string,
	args ...interface{},
) {

	e.logMessages = append([]string{fmt.Sprintf(format, args...)}, e.logMessages...)
}

func (code ErrorCode) String() string {
	return string(code)
}

func IsErrNoRows(err error) bool {

	wrappedError, ok := err.(*Error)
	if !ok {
		return errors.Is(err, sql.ErrNoRows)
	}

	return errors.Is(wrappedError.Err(), sql.ErrNoRows)
}

func (e *Error) HttpStatus() int {

	if e.httpStatusCode != 0 {
		return e.httpStatusCode
	}

	statusCode, ok := httpStatusErrorCodeMap[e.errorCode]
	if !ok || statusCode == 0 {
		return http.StatusBadRequest
	}

	return statusCode
}

func (e *Error) LogErrorMessages() {
	switch e.notify {
	case true:
		logger.Error(e.Error())
	default:
		logger.Warn(e.Error())
	}
}

func (e *Error) Notify() *Error {
	e.notify = true
	return e
}

package apperror

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

type AppError struct {
	Code    int
	Message string
	err     error
	stack   []byte
}

var includeStackTrace = false

func SetIncludeStackTrace(v bool) {
	includeStackTrace = v
}

func NewAppError(code int, message string, err error) *AppError {
	if includeStackTrace {
		return NewAppErrorWithTrace(code, message, err)
	}

	return &AppError{
		Code:    code,
		Message: message,
		err:     err,
	}
}

func NewAppErrorWithTrace(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		err:     err,
		stack:   debug.Stack(),
	}
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("(%d) %s: %s", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("(%d) %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) GetStackTrace() []byte {
	return e.stack
}

func (e *AppError) ContainsStackTrace() bool {
	return len(e.stack) > 0
}

func IsErrorCode(err error, code int) bool {
	aerr, ok := err.(*AppError)
	return ok && aerr.Code == code
}

func NewInternal(err error) error {
	return NewAppErrorWithTrace(CodeInternal, "internal error", err)
}

func Wrap(err error) error {
	if _, ok := err.(*AppError); ok {
		return err
	}
	if errors.Is(err, context.Canceled) {
		return NewCanceled(err)
	}
	return NewInternal(err)
}

func NewInternalFmt(format string, args ...interface{}) error {
	return NewInternal(fmt.Errorf(format, args...))
}

func NewTypeAssertionFailed(want interface{}, got interface{}) error {
	return NewInternalFmt("type assert: want %T, got %T", want, got)
}

func NewCanceled(err error) error {
	return NewAppError(CodeCanceled, "canceled", err)
}

func NewBadRequest(err error) error {
	return NewAppError(
		CodeBadRequest,
		"bad request",
		err,
	)
}

func NewValidationFailed(err error) error {
	return NewAppError(
		CodeValidationFailed,
		"validation failed",
		err,
	)
}

func NewConstraintViolation(err error) error {
	return NewAppError(
		CodeConstraintViolation,
		"validation failed",
		err,
	)
}

func NewNotFound() error {
	return NewAppError(
		CodeNotFound,
		"not found",
		nil,
	)
}

func NewEntityNotFound(name string) error {
	return NewAppError(
		CodeNotFound,
		fmt.Sprintf("%s not found", name),
		nil,
	)
}

func NewAlreadyExists(name string) error {
	return NewAppError(
		CodeAlreadyExists,
		fmt.Sprintf("%s already exists", name),
		nil,
	)
}

func NewUnauthorized(err error) error {
	return NewAppError(
		CodeUnauthorized,
		"unauthorized",
		err,
	)
}

func NewWrongPassword(err error) error {
	return NewAppError(
		CodeUnauthorized,
		"wrong password",
		err,
	)
}

func NewInvalidToken(err error) error {
	return NewAppError(
		CodeUnauthorized,
		"invalid token",
		err,
	)
}

func NewForbidden(err error) error {
	return NewAppError(
		CodeForbidden,
		"forbidden",
		err,
	)
}

func NewImageSizeExceeded(maxSize string) error {
	return NewAppError(
		CodeBadRequest,
		fmt.Sprintf("image file size must be less than %s", maxSize),
		nil,
	)
}

func NewRestrictredFileType(types ...string) error {
	sb := strings.Builder{}
	sb.WriteString("file type must be ")
	for i := 0; i < len(types); i++ {
		if i == len(types)-1 {
			sb.WriteString(fmt.Sprintf("or %s", types[i]))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s, ", types[i]))
	}

	return NewAppError(
		CodeBadRequest,
		sb.String(),
		nil,
	)
}

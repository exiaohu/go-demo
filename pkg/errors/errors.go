package errors

import (
	"errors"
	"fmt"
)

// ErrorType 错误类型

type ErrorType int

const (
	// ErrTypeUnknown 未知错误
	ErrTypeUnknown ErrorType = iota
	// ErrTypeValidation 验证错误
	ErrTypeValidation
	// ErrTypeNotFound 未找到错误
	ErrTypeNotFound
	// ErrTypeUnauthorized 未授权错误
	ErrTypeUnauthorized
	// ErrTypeForbidden 禁止访问错误
	ErrTypeForbidden
	// ErrTypeInternal 内部服务器错误
	ErrTypeInternal
)

// AppError 自定义应用程序错误

type AppError struct {
	Type    ErrorType `json:"type"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// New 创建新的应用程序错误
func New(errType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errType,
		Code:    errorCode(errType),
		Message: message,
	}
}

// NewWithDetails 创建包含详细信息的应用程序错误
func NewWithDetails(errType ErrorType, message, details string) *AppError {
	return &AppError{
		Type:    errType,
		Code:    errorCode(errType),
		Message: message,
		Details: details,
	}
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// errorCode 根据错误类型返回 HTTP 状态码
func errorCode(errType ErrorType) int {
	switch errType {
	case ErrTypeValidation:
		return 400
	case ErrTypeNotFound:
		return 404
	case ErrTypeUnauthorized:
		return 401
	case ErrTypeForbidden:
		return 403
	case ErrTypeInternal:
		return 500
	case ErrTypeUnknown:
		return 500
	default:
		return 500
	}
}

// String implements error interface
func (e *AppError) String() string {
	return e.Error()
}

// Is reports whether err matches target
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// IsType checks if error is of specific type
func IsType(err error, errType ErrorType) bool {
	var e *AppError
	if errors.As(err, &e) {
		return e.Type == errType
	}
	return false
}

// GetType returns the error type
func GetType(err error) ErrorType {
	var e *AppError
	if errors.As(err, &e) {
		return e.Type
	}
	return ErrTypeUnknown
}

// GetDetails returns error details
func GetDetails(err error) string {
	var e *AppError
	if errors.As(err, &e) {
		return e.Details
	}
	return ""
}

// IsValidationError checks if error is validation error
func IsValidationError(err error) bool {
	return IsType(err, ErrTypeValidation)
}

// IsNotFoundError checks if error is not found error
func IsNotFoundError(err error) bool {
	return IsType(err, ErrTypeNotFound)
}

// IsUnauthorizedError 检查是否是未授权错误
func IsUnauthorizedError(err error) bool {
	return IsType(err, ErrTypeUnauthorized)
}

// IsForbiddenError 检查是否是禁止访问错误
func IsForbiddenError(err error) bool {
	return IsType(err, ErrTypeForbidden)
}

// IsInternalError 检查是否是内部服务器错误
func IsInternalError(err error) bool {
	return IsType(err, ErrTypeInternal)
}

package math

import (
	"github.com/exiaohu/go-demo/pkg/errors"
)

// Add 加法函数
func Add(a, b int) (int, error) {
	return a + b, nil
}

// Subtract 减法函数
func Subtract(a, b int) (int, error) {
	return a - b, nil
}

// Multiply 乘法函数
func Multiply(a, b int) (int, error) {
	return a * b, nil
}

// Divide 除法函数
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New(errors.ErrTypeValidation, "Division by zero")
	}
	return a / b, nil
}

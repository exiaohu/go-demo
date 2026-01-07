package handler

import (
	"net/http"
	"strconv"

	"github.com/exiaohu/go-demo/internal/math"
	"github.com/exiaohu/go-demo/pkg/errors"
	"github.com/exiaohu/go-demo/pkg/response"
)

// HomeHandler 主页
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Playground!"))
}

// HealthCheckHandler 健康检查
func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// AddHandler 加法运算
// @Summary Add two integers
// @Description get sum of two integers
// @Tags math
// @Accept  json
// @Produce  json
// @Param a query int true "First integer"
// @Param b query int true "Second integer"
// @Success 200 {string} string "Result"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /add [get]
func AddHandler(w http.ResponseWriter, r *http.Request) {
	handleMathRequest(w, r, math.Add)
}

// SubtractHandler 减法运算
// @Summary Subtract two integers
// @Description get difference of two integers
// @Tags math
// @Accept  json
// @Produce  json
// @Param a query int true "First integer"
// @Param b query int true "Second integer"
// @Success 200 {string} string "Result"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /subtract [get]
func SubtractHandler(w http.ResponseWriter, r *http.Request) {
	handleMathRequest(w, r, math.Subtract)
}

// MultiplyHandler 乘法运算
// @Summary Multiply two integers
// @Description get product of two integers
// @Tags math
// @Accept  json
// @Produce  json
// @Param a query int true "First integer"
// @Param b query int true "Second integer"
// @Success 200 {string} string "Result"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /multiply [get]
func MultiplyHandler(w http.ResponseWriter, r *http.Request) {
	handleMathRequest(w, r, math.Multiply)
}

// DivideHandler 除法运算
// @Summary Divide two integers
// @Description get quotient of two integers
// @Tags math
// @Accept  json
// @Produce  json
// @Param a query int true "First integer"
// @Param b query int true "Second integer"
// @Success 200 {string} string "Result"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /divide [get]
func DivideHandler(w http.ResponseWriter, r *http.Request) {
	handleMathRequest(w, r, math.Divide)
}

func handleMathRequest(w http.ResponseWriter, r *http.Request, op func(int, int) (int, error)) {
	if r.Method != http.MethodGet {
		response.Error(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	a, err := parseIntParam(query.Get("a"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	b, err := parseIntParam(query.Get("b"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	result, err := op(a, b)
	if err != nil {
		if errors.IsValidationError(err) {
			response.Error(w, r, http.StatusBadRequest, err.Error())
		} else {
			response.Error(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(w, r, map[string]int{"result": result})
}

func parseIntParam(param string) (int, error) {
	if param == "" {
		return 0, errors.New(errors.ErrTypeValidation, "Parameter is required")
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0, errors.NewWithDetails(errors.ErrTypeValidation, "Invalid parameter format", err.Error())
	}
	return val, nil
}

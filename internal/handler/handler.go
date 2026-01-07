package handler

import (
	"context"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/exiaohu/go-demo/internal/service"
	"github.com/exiaohu/go-demo/pkg/errors"
	"github.com/exiaohu/go-demo/pkg/response"
	"github.com/exiaohu/go-demo/pkg/util/ip"
)

type Handler struct {
	calcService service.CalculatorService
}

func NewHandler(calcService service.CalculatorService) *Handler {
	return &Handler{
		calcService: calcService,
	}
}

// HomeHandler 主页
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Playground!"))
}

// HealthCheckHandler 健康检查
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
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
func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	h.handleMathRequest(w, r, h.calcService.Add)
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
func (h *Handler) SubtractHandler(w http.ResponseWriter, r *http.Request) {
	h.handleMathRequest(w, r, h.calcService.Subtract)
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
func (h *Handler) MultiplyHandler(w http.ResponseWriter, r *http.Request) {
	h.handleMathRequest(w, r, h.calcService.Multiply)
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
func (h *Handler) DivideHandler(w http.ResponseWriter, r *http.Request) {
	h.handleMathRequest(w, r, h.calcService.Divide)
}

func (h *Handler) handleMathRequest(w http.ResponseWriter, r *http.Request, op func(context.Context, int, int, string) (int, error)) {
	tr := otel.Tracer("handler")
	ctx, span := tr.Start(r.Context(), "handleMathRequest")
	defer span.End()

	if r.Method != http.MethodGet {
		span.RecordError(errors.New(errors.ErrTypeValidation, "Method not allowed"))
		response.Error(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	a, err := parseIntParam(query.Get("a"))
	if err != nil {
		span.RecordError(err)
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	b, err := parseIntParam(query.Get("b"))
	if err != nil {
		span.RecordError(err)
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	span.SetAttributes(
		attribute.Int("a", a),
		attribute.Int("b", b),
	)

	result, err := op(ctx, a, b, ip.GetClientIP(r))
	if err != nil {
		span.RecordError(err)
		if errors.IsValidationError(err) {
			response.Error(w, r, http.StatusBadRequest, err.Error())
		} else {
			response.Error(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	span.SetAttributes(attribute.Int("result", result))
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

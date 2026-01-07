package handler

import (
	"net/http"
	"strconv"

	"github.com/exiaohu/go-demo/pkg/response"
)

// HistoryHandler 获取计算历史
// @Summary Get calculation history
// @Description get latest calculation history
// @Tags history
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit number of records"
// @Success 200 {object} response.Response{data=[]model.CalculationHistory} "History records"
// @Failure 500 {string} string "Internal Server Error"
// @Router /history [get]
func (h *Handler) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := h.calcService.GetHistory(r.Context(), limit)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "Failed to fetch history")
		return
	}

	response.Success(w, r, history)
}

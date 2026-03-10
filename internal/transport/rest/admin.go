package rest

import (
	"net/http"

	"financial_system/internal/domain"
)

// getAllLogs godoc
// @Summary      Все логи действий
// @Description  Возвращает все записи action_logs в порядке убывания времени.
// @Tags         admin
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   domain.ActionLog
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/logs [get]
func (h *Handler) getAllLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.services.Admin.GetAllLogs()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось получить логи")
		return
	}
	if logs == nil {
		logs = []domain.ActionLog{}
	}
	respondJSON(w, http.StatusOK, logs)
}

// undoAction godoc
// @Summary      Отменить действие
// @Tags         admin
// @Security     BearerAuth
// @Param        id   path  int  true  "ID записи лога"
// @Router       /admin/logs/{id}/undo [post]
func (h *Handler) undoAction(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNotImplemented, map[string]string{"message": "undo пока не реализован"})
}

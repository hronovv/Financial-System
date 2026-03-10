package rest

import (
	"net/http"
	"strconv"

	"financial_system/internal/domain"

	"github.com/gorilla/mux"
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
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/logs/{id}/undo [post]
func (h *Handler) undoAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		respondError(w, http.StatusBadRequest, "id лога обязателен")
		return
	}
	logID, err := strconv.Atoi(idStr)
	if err != nil || logID <= 0 {
		respondError(w, http.StatusBadRequest, "id лога должен быть положительным числом")
		return
	}

	if err := h.services.Admin.UndoAction(logID, h.services.Repositories); err != nil {
		switch err {
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "запись лога не найдена")
		default:
			respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package rest

import (
	"net/http"
	"strconv"

	"financial_system/internal/domain"

	"github.com/gorilla/mux"
)

// getAllLogs godoc
// @Summary      All action logs
// @Description  Returns all action_logs ordered by created_at desc.
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
// @Summary      Undo action
// @Description  Logical undo of the action by log entry. Not for auth_sign_up/auth_sign_in.
// @Tags         admin
// @Security     BearerAuth
// @Param        id   path  int  true  "Log entry ID"
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
		case domain.ErrActionAlreadyUndone:
			respondError(w, http.StatusBadRequest, "действие уже было отменено")
		default:
			respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// undoAllActions godoc
// @Summary      Undo all actions
// @Description  Undoes every undoable client/manager action (newest first). Skips auth_sign_up, auth_sign_in and already undone.
// @Tags         admin
// @Security     BearerAuth
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/logs/undo-all [post]
func (h *Handler) undoAllActions(w http.ResponseWriter, r *http.Request) {
	if err := h.services.Admin.UndoAllActions(h.services.Repositories); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

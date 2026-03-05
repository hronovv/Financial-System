package rest

import "net/http"

// getAllLogs godoc
// @Summary      Все логи действий
// @Tags         admin
// @Security     BearerAuth
// @Router       /admin/logs [get]
func (h *Handler) getAllLogs(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// undoAction godoc
// @Summary      Отменить действие
// @Tags         admin
// @Security     BearerAuth
// @Param        id   path  int  true  "ID записи лога"
// @Router       /admin/logs/{id}/undo [post]
func (h *Handler) undoAction(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

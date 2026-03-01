package rest

import "net/http"

func (h *Handler) getAllLogs(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) undoAction(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

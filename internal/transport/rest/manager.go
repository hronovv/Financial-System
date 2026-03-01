package rest

import "net/http"

func (h *Handler) approveUser(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) getUserHistory(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) blockAccount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) unblockAccount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) getEnterprisesWithEmployees(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}


func (h *Handler) addEmployeeToEnterprise(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) removeEmployeeFromEnterprise(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) approveSalaryApplication(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}
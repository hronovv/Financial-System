package rest

import "net/http"

func (h *Handler) getBanks(w http.ResponseWriter, r *http.Request) {
	banks, err := h.services.Bank.GetBanks()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Не удалось получить список банков")
		return
	}

	respondJSON(w, http.StatusOK, banks)
}

func (h *Handler) getEnterprises(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) openAccount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) closeAccount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) transferFromAccount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) getAccountHistory(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) openDeposit(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) closeDeposit(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}	

func (h *Handler) transferFromDeposit(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

func (h *Handler) accumulateDeposit(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}	

func (h *Handler) applyForSalaryProject(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}
func (h *Handler) receiveSalary(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

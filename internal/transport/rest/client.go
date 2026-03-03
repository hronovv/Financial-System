package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"financial_system/internal/domain"

	"github.com/gorilla/mux"
)

func (h *Handler) getBanks(w http.ResponseWriter, r *http.Request) {
	banks, err := h.services.Bank.GetBanks()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось получить список банков")
		return
	}

	respondJSON(w, http.StatusOK, banks)
}

func (h *Handler) getEnterprises(w http.ResponseWriter, r *http.Request) {
	enterprises, err := h.services.Enterprise.GetEnterprises()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось получить список предприятий")
		return
	}

	respondJSON(w, http.StatusOK, enterprises)
}

func (h *Handler) openAccount(w http.ResponseWriter, r *http.Request) {
	var input struct {
		// TODO: user_id в будущем нужно будет получать из JWT, а не из тела запроса.
		UserID int `json:"user_id"`
		BankID int `json:"bank_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.UserID <= 0 {
		respondError(w, http.StatusBadRequest, "user_id должен быть положительным числом")
		return
	}
	if input.BankID <= 0 {
		respondError(w, http.StatusBadRequest, "bank_id должен быть положительным числом")
		return
	}

	account, err := h.services.Account.OpenAccount(input.UserID, input.BankID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось открыть счет")
		return
	}

	respondJSON(w, http.StatusCreated, account)
}

func (h *Handler) closeAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		respondError(w, http.StatusBadRequest, "id счета обязателен")
		return
	}

	accountID, err := strconv.Atoi(idStr)
	if err != nil || accountID <= 0 {
		respondError(w, http.StatusBadRequest, "id счета должен быть положительным числом")
		return
	}

	var input struct {
		// TODO: user_id в будущем нужно будет получать из JWT, а не из тела запроса.
		UserID int `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.UserID <= 0 {
		respondError(w, http.StatusBadRequest, "user_id должен быть положительным числом")
		return
	}

	err = h.services.Account.CloseAccount(input.UserID, accountID)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			respondError(w, http.StatusForbidden, "недостаточно прав для закрытия этого счета")
		case domain.ErrAccountAlreadyClosed:
			respondError(w, http.StatusBadRequest, "счет уже закрыт")
		case domain.ErrAccountHasNonZeroBalance:
			respondError(w, http.StatusBadRequest, "нельзя закрыть счет с ненулевым балансом")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось закрыть счет")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

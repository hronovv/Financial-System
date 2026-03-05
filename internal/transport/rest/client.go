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
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input struct {
		BankID int `json:"bank_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.BankID <= 0 {
		respondError(w, http.StatusBadRequest, "bank_id должен быть положительным числом")
		return
	}

	account, err := h.services.Account.OpenAccount(userID, input.BankID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось открыть счет")
		return
	}

	respondJSON(w, http.StatusCreated, account)
}

func (h *Handler) closeAccount(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

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

	err = h.services.Account.CloseAccount(userID, accountID)
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
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input struct {
		FromAccountID int     `json:"from_account_id"`
		ToAccountID   *int    `json:"to_account_id,omitempty"`
		ToDepositID   *int    `json:"to_deposit_id,omitempty"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.FromAccountID <= 0 {
		respondError(w, http.StatusBadRequest, "from_account_id должен быть положительным числом")
		return
	}
	if input.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount должен быть больше 0")
		return
	}

	err := h.services.Account.TransferFromAccount(userID, input.FromAccountID, input.ToAccountID, input.ToDepositID, input.Amount)
	if err != nil {
		switch err {
		case domain.ErrInvalidAmount, domain.ErrInvalidTransferTarget:
			respondError(w, http.StatusBadRequest, "неверные параметры перевода")
		case domain.ErrForbidden:
			respondError(w, http.StatusForbidden, "недостаточно прав")
		case domain.ErrInsufficientFunds:
			respondError(w, http.StatusBadRequest, "недостаточно средств")
		case domain.ErrAccountBlocked, domain.ErrDepositBlocked:
			respondError(w, http.StatusBadRequest, "счет/вклад заблокирован")
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "счет/вклад не найден")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось выполнить перевод")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getAccountHistory(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		respondError(w, http.StatusBadRequest, "параметр account_id обязателен")
		return
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil || accountID <= 0 {
		respondError(w, http.StatusBadRequest, "account_id должен быть положительным числом")
		return
	}

	history, err := h.services.Account.GetAccountHistory(userID, accountID)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			respondError(w, http.StatusForbidden, "недостаточно прав для просмотра истории этого счета")
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "счет не найден")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось получить историю счета")
		}
		return
	}

	respondJSON(w, http.StatusOK, history)
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

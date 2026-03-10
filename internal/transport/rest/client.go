package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"financial_system/internal/domain"

	"github.com/gorilla/mux"
)

// openAccountRequest is the body for POST /client/accounts.
type openAccountRequest struct {
	BankID int `json:"bank_id" example:"1"`
}

// openDepositRequest is the body for POST /client/deposits.
type openDepositRequest struct {
	BankID       int     `json:"bank_id" example:"1"`
	InterestRate float64 `json:"interest_rate" example:"5.5"`
}

// transferFromAccountRequest is the body for POST /client/accounts/transfer.
type transferFromAccountRequest struct {
	FromAccountID int      `json:"from_account_id" example:"10"`
	ToAccountID   *int     `json:"to_account_id,omitempty" example:"11"`
	ToDepositID   *int     `json:"to_deposit_id,omitempty" example:"5"`
	Amount        float64 `json:"amount" example:"100.50"`
}

// applySalaryProjectRequest is the body for POST /client/salary-project/apply.
type applySalaryProjectRequest struct {
	EnterpriseID int     `json:"enterprise_id" example:"1"`
	Amount       float64 `json:"amount" example:"50000"`
}

// transferFromDepositRequest is the body for POST /client/deposits/transfer.
type transferFromDepositRequest struct {
	FromDepositID int   `json:"from_deposit_id" example:"3"`
	ToAccountID   *int  `json:"to_account_id,omitempty" example:"10"`
	ToDepositID   *int  `json:"to_deposit_id,omitempty" example:"5"`
	Amount        float64 `json:"amount" example:"100.50"`
}

// accumulateDepositRequest is the body for POST /client/deposits/{id}/accumulate.
type accumulateDepositRequest struct {
	FromAccountID int     `json:"from_account_id" example:"10"`
	Amount        float64 `json:"amount" example:"500"`
}

// receiveSalaryRequest is the body for POST /client/salary-project/receive.
type receiveSalaryRequest struct {
	ApplicationID int  `json:"application_id" example:"1"`
	ToAccountID   *int `json:"to_account_id,omitempty" example:"5"`
	ToDepositID   *int `json:"to_deposit_id,omitempty" example:"2"`
}

// getBanks godoc
// @Summary      List banks
// @Description  Returns all banks. Requires client role and JWT.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.Bank
// @Failure      401  {object}  map[string]string
// @Router       /client/banks [get]
func (h *Handler) getBanks(w http.ResponseWriter, r *http.Request) {
	banks, err := h.services.Bank.GetBanks()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось получить список банков")
		return
	}

	respondJSON(w, http.StatusOK, banks)
}

// getEnterprises godoc
// @Summary      List enterprises
// @Description  Returns all enterprises. Requires client role and JWT.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.Enterprise
// @Failure      401  {object}  map[string]string
// @Router       /client/enterprises [get]
func (h *Handler) getEnterprises(w http.ResponseWriter, r *http.Request) {
	enterprises, err := h.services.Enterprise.GetEnterprises()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось получить список предприятий")
		return
	}

	respondJSON(w, http.StatusOK, enterprises)
}

// openAccount godoc
// @Summary      Open account
// @Description  Opens an account in the given bank. user_id from JWT.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  openAccountRequest  true  "bank_id"
// @Success      201  {object}  domain.Account
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /client/accounts [post]
func (h *Handler) openAccount(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input openAccountRequest

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
		respondError(w, http.StatusBadRequest, "не удалось открыть счет")
		return
	}

	uid := userID
	_ = h.services.Audit.LogAction(&uid, "client_open_account", map[string]any{
		"account_id":     account.ID,
		"bank_id":        input.BankID,
		"account_number": account.AccountNumber,
	})

	respondJSON(w, http.StatusCreated, account)
}

// closeAccount godoc
// @Summary      Close account
// @Description  Closes account (is_blocked=true). Balance must be 0. user_id from JWT.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Account ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /client/accounts/{id} [delete]
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

	uid := userID
	_ = h.services.Audit.LogAction(&uid, "client_close_account", map[string]any{
		"account_id": accountID,
	})

	w.WriteHeader(http.StatusNoContent)
}

// transferFromAccount godoc
// @Summary      Transfer from account
// @Description  Transfer from account to another account or deposit (same user). Provide either to_account_id or to_deposit_id.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  transferFromAccountRequest  true  "from_account_id, to_account_id or to_deposit_id, amount"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /client/accounts/transfer [post]
func (h *Handler) transferFromAccount(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input transferFromAccountRequest

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

	uid := userID
	details := map[string]any{
		"from_account_id": input.FromAccountID,
		"amount":          input.Amount,
	}
	if input.ToAccountID != nil {
		details["to_account_id"] = *input.ToAccountID
	}
	if input.ToDepositID != nil {
		details["to_deposit_id"] = *input.ToDepositID
	}
	_ = h.services.Audit.LogAction(&uid, "client_transfer_from_account", details)

	w.WriteHeader(http.StatusNoContent)
}

// getAccountHistory godoc
// @Summary      Account history
// @Description  Returns transactions for the account. user_id from JWT. Query: account_id.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account_id  query     int  true  "Account ID"
// @Success      200  {array}   domain.Transaction
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /client/accounts/history [get]
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

// openDeposit godoc
// @Summary      Open deposit
// @Description  Creates a deposit in a bank with given interest rate. Initial balance 0.
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  openDepositRequest  true  "bank_id, interest_rate"
// @Success      201  {object}  domain.Deposit
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /client/deposits [post]
func (h *Handler) openDeposit(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input openDepositRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.BankID <= 0 {
		respondError(w, http.StatusBadRequest, "bank_id должен быть положительным числом")
		return
	}
	if input.InterestRate < 0 {
		respondError(w, http.StatusBadRequest, "interest_rate не может быть отрицательным")
		return
	}

	deposit, err := h.services.Deposit.OpenDeposit(userID, input.BankID, input.InterestRate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось открыть вклад")
		return
	}

	uid := userID
	_ = h.services.Audit.LogAction(&uid, "client_open_deposit", map[string]any{
		"deposit_id":    deposit.ID,
		"bank_id":       input.BankID,
		"interest_rate": input.InterestRate,
	})

	respondJSON(w, http.StatusCreated, deposit)
}

// closeDeposit godoc
// @Summary      Close deposit
// @Description  Closes deposit (is_blocked=true). Balance must be 0.
// @Tags         client
// @Security     BearerAuth
// @Param        id   path      int  true  "Deposit ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /client/deposits/{id} [delete]
func (h *Handler) closeDeposit(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		respondError(w, http.StatusBadRequest, "id вклада обязателен")
		return
	}

	depositID, err := strconv.Atoi(idStr)
	if err != nil || depositID <= 0 {
		respondError(w, http.StatusBadRequest, "id вклада должен быть положительным числом")
		return
	}

	err = h.services.Deposit.CloseDeposit(userID, depositID)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			respondError(w, http.StatusForbidden, "недостаточно прав для закрытия этого вклада")
		case domain.ErrDepositAlreadyClosed:
			respondError(w, http.StatusBadRequest, "вклад уже закрыт")
		case domain.ErrDepositHasNonZeroBalance:
			respondError(w, http.StatusBadRequest, "нельзя закрыть вклад с ненулевым балансом")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось закрыть вклад")
		}
		return
	}

	uid := userID
	_ = h.services.Audit.LogAction(&uid, "client_close_deposit", map[string]any{
		"deposit_id": depositID,
	})

	w.WriteHeader(http.StatusNoContent)
}

// transferFromDeposit godoc
// @Summary      Transfer from deposit
// @Description  Transfer from deposit to account or another deposit (same user). Provide exactly one of to_account_id or to_deposit_id.
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Param        body  body  transferFromDepositRequest  true  "from_deposit_id, to_account_id or to_deposit_id, amount"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /client/deposits/transfer [post]
func (h *Handler) transferFromDeposit(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input transferFromDepositRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.FromDepositID <= 0 {
		respondError(w, http.StatusBadRequest, "from_deposit_id должен быть положительным числом")
		return
	}
	if input.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount должен быть больше 0")
		return
	}
	hasAccount := input.ToAccountID != nil && *input.ToAccountID > 0
	hasDeposit := input.ToDepositID != nil && *input.ToDepositID > 0
	if !hasAccount && !hasDeposit {
		respondError(w, http.StatusBadRequest, "укажите to_account_id или to_deposit_id")
		return
	}
	if hasAccount && hasDeposit {
		respondError(w, http.StatusBadRequest, "укажите только один из to_account_id или to_deposit_id")
		return
	}

	err := h.services.Deposit.TransferFromDeposit(userID, input.FromDepositID, input.ToAccountID, input.ToDepositID, input.Amount)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "вклад не найден")
		case domain.ErrAccountNotFound:
			respondError(w, http.StatusNotFound, "счёт не найден")
		case domain.ErrDepositNotFound:
			respondError(w, http.StatusNotFound, "вклад (приёмник) не найден")
		case domain.ErrForbidden:
			respondError(w, http.StatusForbidden, "недостаточно прав")
		case domain.ErrInvalidAmount, domain.ErrInvalidTransferTarget:
			respondError(w, http.StatusBadRequest, "неверные параметры перевода")
		case domain.ErrInsufficientFunds:
			respondError(w, http.StatusBadRequest, "недостаточно средств на вкладе")
		case domain.ErrDepositBlocked:
			respondError(w, http.StatusBadRequest, "вклад заблокирован")
		case domain.ErrAccountBlocked:
			respondError(w, http.StatusBadRequest, "счёт заблокирован")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось выполнить перевод")
		}
		return
	}

	uid := userID
	details := map[string]any{
		"from_deposit_id": input.FromDepositID,
		"amount":          input.Amount,
	}
	if input.ToAccountID != nil {
		details["to_account_id"] = *input.ToAccountID
	}
	if input.ToDepositID != nil {
		details["to_deposit_id"] = *input.ToDepositID
	}
	_ = h.services.Audit.LogAction(&uid, "client_transfer_from_deposit", details)

	w.WriteHeader(http.StatusNoContent)
}

// accumulateDeposit godoc
// @Summary      Accumulate deposit (top-up from account)
// @Description  Transfers from user account to the deposit (id in path).
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  int  true  "Deposit ID (target)"
// @Param        body  body  accumulateDepositRequest  true  "from_account_id, amount"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /client/deposits/{id}/accumulate [post]
func (h *Handler) accumulateDeposit(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	vars := mux.Vars(r)
	depositIDStr, ok := vars["id"]
	if !ok {
		respondError(w, http.StatusBadRequest, "id вклада обязателен")
		return
	}
	depositID, err := strconv.Atoi(depositIDStr)
	if err != nil || depositID <= 0 {
		respondError(w, http.StatusBadRequest, "id вклада должен быть положительным числом")
		return
	}

	var input accumulateDepositRequest
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

	err = h.services.Account.TransferFromAccount(userID, input.FromAccountID, nil, &depositID, input.Amount)
	if err != nil {
		switch err {
		case domain.ErrInvalidAmount, domain.ErrInvalidTransferTarget:
			respondError(w, http.StatusBadRequest, "неверные параметры перевода")
		case domain.ErrForbidden:
			respondError(w, http.StatusForbidden, "недостаточно прав")
		case domain.ErrInsufficientFunds:
			respondError(w, http.StatusBadRequest, "недостаточно средств на счёте")
		case domain.ErrAccountBlocked, domain.ErrDepositBlocked:
			respondError(w, http.StatusBadRequest, "счёт или вклад заблокирован")
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "счёт или вклад не найден")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось пополнить вклад")
		}
		return
	}

	uid := userID
	_ = h.services.Audit.LogAction(&uid, "client_accumulate_deposit", map[string]any{
		"from_account_id": input.FromAccountID,
		"deposit_id":      depositID,
		"amount":          input.Amount,
	})

	w.WriteHeader(http.StatusNoContent)
}

// applyForSalaryProject godoc
// @Summary      Apply for salary project
// @Description  Creates application with status pending. Only enterprise employee.
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Param        body  body  applySalaryProjectRequest  true  "enterprise_id, amount"
// @Success      201  {object}  domain.SalaryApplication
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /client/salary-project/apply [post]
func (h *Handler) applyForSalaryProject(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input applySalaryProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.EnterpriseID <= 0 {
		respondError(w, http.StatusBadRequest, "enterprise_id должен быть положительным числом")
		return
	}
	if input.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount должен быть больше 0")
		return
	}

	app, err := h.services.SalaryProject.ApplyForSalaryProject(userID, input.EnterpriseID, input.Amount)
	if err != nil {
		switch err {
		case domain.ErrNotEmployee:
			respondError(w, http.StatusForbidden, "вы не являетесь сотрудником этого предприятия")
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "предприятие не найдено")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось подать заявку")
		}
		return
	}

	uid := userID
	_ = h.services.Audit.LogAction(&uid, "client_salary_application_create", map[string]any{
		"application_id": app.ID,
		"enterprise_id":  input.EnterpriseID,
		"amount":         input.Amount,
	})

	respondJSON(w, http.StatusCreated, app)
}

// receiveSalary godoc
// @Summary      Receive salary
// @Description  Credits salary for approved application to given account or deposit (one of to_account_id, to_deposit_id required).
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Param        body  body  receiveSalaryRequest  true  "application_id, to_account_id or to_deposit_id"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /client/salary-project/receive [post]
func (h *Handler) receiveSalary(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var input receiveSalaryRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.ApplicationID <= 0 {
		respondError(w, http.StatusBadRequest, "application_id должен быть положительным числом")
		return
	}
	hasAccount := input.ToAccountID != nil && *input.ToAccountID > 0
	hasDeposit := input.ToDepositID != nil && *input.ToDepositID > 0
	if !hasAccount && !hasDeposit {
		respondError(w, http.StatusBadRequest, "укажите to_account_id или to_deposit_id")
		return
	}
	if hasAccount && hasDeposit {
		respondError(w, http.StatusBadRequest, "укажите только один из to_account_id или to_deposit_id")
		return
	}

	err := h.services.SalaryProject.ReceiveSalary(userID, input.ApplicationID, input.ToAccountID, input.ToDepositID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "заявка не найдена")
		case domain.ErrAccountNotFound:
			respondError(w, http.StatusNotFound, "счёт не найден")
		case domain.ErrDepositNotFound:
			respondError(w, http.StatusNotFound, "вклад не найден")
		case domain.ErrForbidden:
			respondError(w, http.StatusForbidden, "недостаточно прав или счёт/вклад не принадлежит вам")
		case domain.ErrApplicationNotApproved:
			respondError(w, http.StatusBadRequest, "заявка не одобрена или уже выплачена")
		case domain.ErrApplicationAlreadyPaid:
			respondError(w, http.StatusBadRequest, "зарплата по этой заявке уже получена")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось получить зарплату")
		}
		return
	}

	uid := userID
	details := map[string]any{
		"application_id": input.ApplicationID,
	}
	if input.ToAccountID != nil {
		details["to_account_id"] = *input.ToAccountID
	}
	if input.ToDepositID != nil {
		details["to_deposit_id"] = *input.ToDepositID
	}
	_ = h.services.Audit.LogAction(&uid, "client_salary_receive", details)

	w.WriteHeader(http.StatusNoContent)
}

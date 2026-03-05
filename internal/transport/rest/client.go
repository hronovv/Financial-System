package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"financial_system/internal/domain"

	"github.com/gorilla/mux"
)

// openAccountRequest тело запроса POST /client/accounts
type openAccountRequest struct {
	BankID int `json:"bank_id" example:"1"`
}

// openDepositRequest тело запроса POST /client/deposits
type openDepositRequest struct {
	BankID       int     `json:"bank_id" example:"1"`
	InterestRate float64 `json:"interest_rate" example:"5.5"`
}

// transferFromAccountRequest тело запроса POST /client/accounts/transfer
type transferFromAccountRequest struct {
	FromAccountID int      `json:"from_account_id" example:"10"`
	ToAccountID   *int     `json:"to_account_id,omitempty" example:"11"`
	ToDepositID   *int     `json:"to_deposit_id,omitempty" example:"5"`
	Amount        float64 `json:"amount" example:"100.50"`
}

// applySalaryProjectRequest тело запроса POST /client/salary-project/apply
type applySalaryProjectRequest struct {
	EnterpriseID int     `json:"enterprise_id" example:"1"`
	Amount       float64 `json:"amount" example:"50000"`
}

// receiveSalaryRequest тело запроса POST /client/salary-project/receive
type receiveSalaryRequest struct {
	ApplicationID int  `json:"application_id" example:"1"`
	ToAccountID   *int `json:"to_account_id,omitempty" example:"5"`
	ToDepositID   *int `json:"to_deposit_id,omitempty" example:"2"`
}

// getBanks godoc
// @Summary      Список банков
// @Description  Возвращает все банки системы. Требуется роль client и JWT.
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
// @Summary      Список предприятий
// @Description  Возвращает все предприятия. Требуется роль client и JWT.
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
// @Summary      Открыть счёт
// @Description  Открывает счёт в указанном банке. user_id берётся из JWT.
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
		respondError(w, http.StatusInternalServerError, "не удалось открыть счет")
		return
	}

	respondJSON(w, http.StatusCreated, account)
}

// closeAccount godoc
// @Summary      Закрыть счёт
// @Description  Закрывает счёт (is_blocked=true). Баланс должен быть 0. user_id из JWT.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID счёта"
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

	w.WriteHeader(http.StatusNoContent)
}

// transferFromAccount godoc
// @Summary      Перевод со счёта
// @Description  Перевод со счёта на другой счёт или вклад (внутри одного пользователя). Указать либо to_account_id, либо to_deposit_id.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  transferFromAccountRequest  true  "from_account_id, to_account_id или to_deposit_id, amount"
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

	w.WriteHeader(http.StatusNoContent)
}

// getAccountHistory godoc
// @Summary      История операций по счёту
// @Description  Возвращает список транзакций по счёту. user_id из JWT. Query: account_id.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account_id  query     int  true  "ID счёта"
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
// @Summary      Открыть вклад
// @Description  Создаёт вклад в банке с указанной процентной ставкой. Баланс при создании 0.
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

	respondJSON(w, http.StatusCreated, deposit)
}

// closeDeposit godoc
// @Summary      Закрыть вклад
// @Description  Закрывает вклад (is_blocked=true). Баланс должен быть 0.
// @Tags         client
// @Security     BearerAuth
// @Param        id   path      int  true  "ID вклада"
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

	w.WriteHeader(http.StatusNoContent)
}

// transferFromDeposit godoc
// @Summary      Перевод со вклада
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Router       /client/deposits/transfer [post]
func (h *Handler) transferFromDeposit(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// accumulateDeposit godoc
// @Summary      Начисление на вклад
// @Tags         client
// @Security     BearerAuth
// @Param        id   path  int  true  "ID вклада"
// @Router       /client/deposits/{id}/accumulate [post]
func (h *Handler) accumulateDeposit(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// applyForSalaryProject godoc
// @Summary      Подать заявку на зарплатный проект
// @Description  Создаёт заявку со статусом pending. Только сотрудник предприятия может подать заявку.
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

	respondJSON(w, http.StatusCreated, app)
}

// receiveSalary godoc
// @Summary      Получить зарплату
// @Description  Зачисляет зарплату по одобренной заявке на указанный счёт или вклад (один из to_account_id, to_deposit_id обязателен).
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Param        body  body  receiveSalaryRequest  true  "application_id, to_account_id или to_deposit_id"
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

	w.WriteHeader(http.StatusNoContent)
}

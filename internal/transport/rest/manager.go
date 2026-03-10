package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"financial_system/internal/domain"

	"github.com/gorilla/mux"
)

// approveUser godoc
// @Summary      Approve client registration
// @Description  Sets is_active = true for user with role client. Client can then sign in.
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "User (client) ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/users/{id}/approve [post]
func (h *Handler) approveUser(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		respondError(w, http.StatusBadRequest, "id пользователя обязателен")
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		respondError(w, http.StatusBadRequest, "id должен быть положительным числом")
		return
	}

	err = h.services.Manager.ApproveUser(userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			respondError(w, http.StatusNotFound, "пользователь не найден")
		case domain.ErrCanOnlyApproveClient:
			respondError(w, http.StatusForbidden, "можно подтверждать только клиентов")
		case domain.ErrUserAlreadyActive:
			respondError(w, http.StatusBadRequest, "пользователь уже подтверждён")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось подтвердить пользователя")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_approve_user", map[string]any{
			"target_user_id": userID,
		})
	}

	w.WriteHeader(http.StatusNoContent)
}

// getUserHistory godoc
// @Summary      User transaction history
// @Description  Merged history for all user accounts, sorted by date.
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "User ID"
// @Success      200  {array}   domain.Transaction
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/users/{id}/history [get]
func (h *Handler) getUserHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		respondError(w, http.StatusBadRequest, "id пользователя обязателен")
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		respondError(w, http.StatusBadRequest, "id должен быть положительным числом")
		return
	}

	history, err := h.services.Manager.GetUserHistory(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось получить историю")
		return
	}

	if history == nil {
		history = []domain.Transaction{}
	}
	respondJSON(w, http.StatusOK, history)
}

// parseAccountIDFromRequest extracts account id from path.
func parseAccountIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return 0, errors.New("id счёта обязателен")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("id счёта должен быть положительным числом")
	}
	return id, nil
}

// blockAccount godoc
// @Summary      Block account
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "Account ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/accounts/{id}/block [post]
func (h *Handler) blockAccount(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	accountID, err := parseAccountIDFromRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Manager.BlockAccount(accountID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "счёт не найден")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось заблокировать счёт")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_block_account", map[string]any{
			"account_id": accountID,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}

// unblockAccount godoc
// @Summary      Unblock account
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "Account ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/accounts/{id}/unblock [post]
func (h *Handler) unblockAccount(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	accountID, err := parseAccountIDFromRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Manager.UnblockAccount(accountID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "счёт не найден")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось разблокировать счёт")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_unblock_account", map[string]any{
			"account_id": accountID,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseDepositIDFromRequest extracts deposit id from path.
func parseDepositIDFromRequest(r *http.Request) (int, error) {
	return parseIDFromPath(r, "id", "id вклада")
}

// blockDeposit godoc
// @Summary      Block deposit
// @Description  Manager can block deposit at any time (no balance check).
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "Deposit ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/deposits/{id}/block [post]
func (h *Handler) blockDeposit(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	depositID, err := parseDepositIDFromRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Manager.BlockDeposit(depositID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "вклад не найден")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось заблокировать вклад")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_block_deposit", map[string]any{
			"deposit_id": depositID,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}

// unblockDeposit godoc
// @Summary      Unblock deposit
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "Deposit ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/deposits/{id}/unblock [post]
func (h *Handler) unblockDeposit(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	depositID, err := parseDepositIDFromRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Manager.UnblockDeposit(depositID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "вклад не найден")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось разблокировать вклад")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_unblock_deposit", map[string]any{
			"deposit_id": depositID,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}

// addEmployeeToEnterpriseRequest is the body for POST /manager/enterprises/{id}/employees.
type addEmployeeToEnterpriseRequest struct {
	UserID int `json:"user_id" example:"3"`
}

// getEnterprisesWithEmployees godoc
// @Summary      Enterprises with employees
// @Tags         manager
// @Security     BearerAuth
// @Success      200  {array}  domain.EnterpriseWithEmployees
// @Failure      401  {object}  map[string]string
// @Router       /manager/enterprises [get]
func (h *Handler) getEnterprisesWithEmployees(w http.ResponseWriter, r *http.Request) {
	list, err := h.services.Manager.GetEnterprisesWithEmployees()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "не удалось получить список предприятий")
		return
	}
	if list == nil {
		list = []domain.EnterpriseWithEmployees{}
	}
	respondJSON(w, http.StatusOK, list)
}

// addEmployeeToEnterprise godoc
// @Summary      Add employee to enterprise
// @Tags         manager
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  int  true  "Enterprise ID"
// @Param        body  body  addEmployeeToEnterpriseRequest  true  "user_id"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/enterprises/{id}/employees [post]
func (h *Handler) addEmployeeToEnterprise(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	enterpriseID, err := parseIDFromPath(r, "id", "id предприятия")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var input addEmployeeToEnterpriseRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.UserID <= 0 {
		respondError(w, http.StatusBadRequest, "user_id должен быть положительным числом")
		return
	}

	err = h.services.Manager.AddEmployee(enterpriseID, input.UserID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "предприятие не найдено")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось добавить сотрудника")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_add_employee", map[string]any{
			"enterprise_id": enterpriseID,
			"user_id":       input.UserID,
		})
	}

	w.WriteHeader(http.StatusNoContent)
}

// removeEmployeeFromEnterprise godoc
// @Summary      Remove employee from enterprise
// @Description  Employee's pending applications for this enterprise are rejected.
// @Tags         manager
// @Security     BearerAuth
// @Param        enterprise_id   path  int  true  "Enterprise ID"
// @Param        user_id        path  int  true  "User ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/enterprises/{enterprise_id}/employees/{user_id} [delete]
func (h *Handler) removeEmployeeFromEnterprise(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	vars := mux.Vars(r)
	enterpriseID, err := strconv.Atoi(vars["enterprise_id"])
	if err != nil || enterpriseID <= 0 {
		respondError(w, http.StatusBadRequest, "enterprise_id должен быть положительным числом")
		return
	}
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil || userID <= 0 {
		respondError(w, http.StatusBadRequest, "user_id должен быть положительным числом")
		return
	}

	err = h.services.Manager.RemoveEmployee(enterpriseID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "предприятие не найдено")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось удалить сотрудника")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_remove_employee", map[string]any{
			"enterprise_id": enterpriseID,
			"user_id":       userID,
		})
	}

	w.WriteHeader(http.StatusNoContent)
}

// parseIDFromPath extracts numeric id from path by key.
func parseIDFromPath(r *http.Request, key, label string) (int, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[key]
	if !ok {
		return 0, errors.New(label + " обязателен")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New(label + " должен быть положительным числом")
	}
	return id, nil
}

// approveSalaryApplication godoc
// @Summary      Approve salary application
// @Description  Approves application (status = approved). Enterprise balance must be >= application amount.
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "Application ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/salary-project/applications/{id}/approve [post]
func (h *Handler) approveSalaryApplication(w http.ResponseWriter, r *http.Request) {
	managerID := userIDFromRequest(r)

	applicationID, err := parseIDFromPath(r, "id", "id заявки")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Manager.ApproveSalaryApplication(applicationID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "заявка не найдена")
		case errors.Is(err, domain.ErrApplicationNotPending):
			respondError(w, http.StatusBadRequest, "заявка уже рассмотрена (не в статусе pending)")
		case errors.Is(err, domain.ErrInsufficientEnterpriseBalance):
			respondError(w, http.StatusBadRequest, "недостаточно средств на балансе предприятия")
		default:
			respondError(w, http.StatusInternalServerError, "не удалось одобрить заявку")
		}
		return
	}

	if managerID != 0 {
		mid := managerID
		_ = h.services.Audit.LogAction(&mid, "manager_approve_salary_application", map[string]any{
			"application_id": applicationID,
		})
	}

	w.WriteHeader(http.StatusNoContent)
}
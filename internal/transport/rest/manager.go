package rest

import (
	"net/http"
	"strconv"

	"financial_system/internal/domain"

	"github.com/gorilla/mux"
)

// approveUser godoc
// @Summary      Подтвердить регистрацию клиента
// @Description  Устанавливает is_active = true для пользователя с ролью client. После этого клиент может войти по SignIn.
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "ID пользователя (клиента)"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /manager/users/{id}/approve [post]
func (h *Handler) approveUser(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusNoContent)
}

// getUserHistory godoc
// @Summary      История операций пользователя
// @Description  Объединённая история по всем счетам пользователя (аналогично /client/accounts/history по каждому счёту), отсортировано по дате.
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "ID пользователя"
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

// blockAccount godoc
// @Summary      Заблокировать счёт
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "ID счёта"
// @Router       /manager/accounts/{id}/block [post]
func (h *Handler) blockAccount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// unblockAccount godoc
// @Summary      Разблокировать счёт
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "ID счёта"
// @Router       /manager/accounts/{id}/unblock [post]
func (h *Handler) unblockAccount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// getEnterprisesWithEmployees godoc
// @Summary      Предприятия с сотрудниками
// @Tags         manager
// @Security     BearerAuth
// @Router       /manager/enterprises [get]
func (h *Handler) getEnterprisesWithEmployees(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// addEmployeeToEnterprise godoc
// @Summary      Добавить сотрудника в предприятие
// @Tags         manager
// @Security     BearerAuth
// @Accept       json
// @Param        id   path  int  true  "ID предприятия"
// @Router       /manager/enterprises/{id}/employees [post]
func (h *Handler) addEmployeeToEnterprise(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// removeEmployeeFromEnterprise godoc
// @Summary      Удалить сотрудника из предприятия
// @Tags         manager
// @Security     BearerAuth
// @Param        enterprise_id   path  int  true  "ID предприятия"
// @Param        user_id         path  int  true  "ID пользователя"
// @Router       /manager/enterprises/{enterprise_id}/employees/{user_id} [delete]
func (h *Handler) removeEmployeeFromEnterprise(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}

// approveSalaryApplication godoc
// @Summary      Подтвердить заявку на зарплатный проект
// @Tags         manager
// @Security     BearerAuth
// @Param        id   path  int  true  "ID заявки"
// @Router       /manager/salary-project/applications/{id}/approve [post]
func (h *Handler) approveSalaryApplication(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, "ok")
}
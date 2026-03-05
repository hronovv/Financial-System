package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"financial_system/internal/domain"
)

// authInputDTO тело запроса для sign-up и sign-in
type authInputDTO struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// signUp godoc
// @Summary      Регистрация клиента
// @Description  Регистрация нового клиента. Требует подтверждения менеджером (is_active).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  authInputDTO  true  "email, password (мин. 8 символов)"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input authInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	input.Email = strings.TrimSpace(input.Email)
	if input.Email == "" {
		respondError(w, http.StatusBadRequest, "email не может быть пустым")
		return
	}
	if len(input.Password) < 8 {
		respondError(w, http.StatusBadRequest, "пароль должен содержать минимум 8 символов")
		return
	}

	err := h.services.Auth.SignUp(input.Email, input.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"message": "Регистрация прошла успешно. Ожидайте подтверждения от менеджера.",
	})
}

// signIn godoc
// @Summary      Вход
// @Description  Вход по email и паролю. Возвращает JWT для заголовка Authorization: Bearer &lt;token&gt;. Аккаунт должен быть подтверждён (is_active).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  authInputDTO  true  "email, password"
// @Success      200  {object}  map[string]string  "token"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input authInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	input.Email = strings.TrimSpace(input.Email)
	if input.Email == "" || input.Password == "" {
		respondError(w, http.StatusBadRequest, "email и пароль обязательны для входа")
		return
	}

	token, err := h.services.Auth.SignIn(input.Email, input.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			respondError(w, http.StatusUnauthorized, "неверный email или пароль")
		case domain.ErrUserNotActive:
			respondError(w, http.StatusForbidden, "аккаунт не подтверждён менеджером")
		default:
			respondError(w, http.StatusInternalServerError, "ошибка входа")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

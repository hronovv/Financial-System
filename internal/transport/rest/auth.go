package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"financial_system/internal/domain"
)

// authInputDTO is the request body for sign-up and sign-in.
type authInputDTO struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// signUp godoc
// @Summary      Client registration
// @Description  Register new client. Requires manager approval (is_active).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  authInputDTO  true  "email, password (min 8 chars)"
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

	_ = h.services.Audit.LogAction(nil, "auth_sign_up", map[string]any{
		"email": input.Email,
	})

	respondJSON(w, http.StatusCreated, map[string]string{
		"message": "Регистрация прошла успешно. Ожидайте подтверждения от менеджера.",
	})
}

// signIn godoc
// @Summary      Sign in
// @Description  Sign in with email and password. Returns JWT for Authorization: Bearer &lt;token&gt;. Account must be approved (is_active).
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

	_ = h.services.Audit.LogAction(nil, "auth_sign_in", map[string]any{
		"email": input.Email,
	})

	respondJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

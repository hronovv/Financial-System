package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"financial_system/internal/domain"
)

type authInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

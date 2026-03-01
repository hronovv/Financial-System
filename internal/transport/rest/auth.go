package rest

import (
	"encoding/json"
	"net/http"
	"strings"
)

type authInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input authInput

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
	if len(input.Password) < 6 {
		respondError(w, http.StatusBadRequest, "пароль должен содержать минимум 6 символов")
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
	var input authInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}
	defer r.Body.Close()

	if input.Email == "" || input.Password == "" {
		respondError(w, http.StatusBadRequest, "email и пароль обязательны для входа")
		return
	}

	// 3. Вызываем Сервис (Бизнес-логику)
	// token, err := h.services.Auth.SignIn(input.Email, input.Password)
	// if err != nil {
	// 	// Если пароль неверный или менеджер еще не подтвердил аккаунт
	// 	respondError(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }

	// 4. Отдаем токен пользователю
	// Пока сервиса нет, отдаем фейковый токен для теста
	token := "fake-jwt-token-123"

	respondJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

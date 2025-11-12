package handlers

import (
	"encoding/json"
	"net/http"
	"taskforge/internal/dto/request"
	"taskforge/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	defer r.Body.Close()

	var req request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		JSONError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	user, err := h.userService.Register(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
	)
	if err != nil {
		switch err {
		case service.ErrEmailAlreadyExist:
			JSONError(w, http.StatusConflict, "Email already registered")
		case service.ErrInvalidEmail:
			JSONError(w, http.StatusBadRequest, "Invalid email format")
		case service.ErrPasswordTooShort:
			JSONError(w, http.StatusBadRequest, "Password too short")
		default:
			JSONError(w, http.StatusInternalServerError, "Registration failed")
		}
		return
	}

	JSONSuccess(w, http.StatusCreated, user)

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	defer r.Body.Close()

	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Email == "" || req.Password == "" {
		JSONError(w, http.StatusBadRequest, "Email and Password are required")
		return
	}

	// Теперь Login возвращает три значения: user, token, error
	user, token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		JSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Возвращаем пользователя и токен
	JSONSuccess(w, http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
		"message": "Login successful",
	})
}	
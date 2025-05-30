package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/service"

	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	// Extract 'username' from query parameters
	username := r.URL.Query().Get("username")
	if username == "" {
		middleware.WriteError(w, http.StatusBadRequest, "Missing 'username' query parameter")
		return
	}

	// Call the service with the username
	resp, err := h.userService.FindByUsername(r.Context(), username)
	if err != nil {
		h.logger.Error("find user by username failed", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Respond with JSON
	middleware.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	data, ok := middleware.GetJSONBody(r)
	if !ok {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Convert the generic data to JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Unmarshal into the request struct
	var req service.GetUserByUserIDRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	resp, err := h.userService.FindByID(r.Context(), req.UserID)
	if err != nil {
		h.logger.Error("find user by ID failed", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	data, ok := middleware.GetJSONBody(r)
	if !ok {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Convert the generic data to JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Unmarshal into the request struct
	var req service.UpdateUsernameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	resp, err := h.userService.UpdateUsername(r.Context(), &req)
	if err != nil {
		h.logger.Error("update username failed", zap.Error(err))
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			middleware.WriteError(w, http.StatusNotFound, "user not found")
		default:
			middleware.WriteError(w, http.StatusInternalServerError, "could not update username")
		}
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

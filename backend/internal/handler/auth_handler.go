package handler

import (
	"encoding/json"
	"net/http"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/service"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
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
	var req service.SignUpRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	resp, err := h.authService.SignUp(r.Context(), req)
	if err != nil {
		h.logger.Error("signup failed", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
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
	var req service.SignInRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	resp, err := h.authService.SignIn(r.Context(), req)
	if err != nil {
		// Log email confirmation errors as INFO
		if err.Error() == "Please check your email for a confirmation link. If you haven't received it, try signing up again." {
			h.logger.Info("signin requires email confirmation",
				zap.String("email", req.Email),
				zap.Error(err),
			)
		} else {
			h.logger.Error("signin failed", zap.Error(err))
		}
		middleware.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

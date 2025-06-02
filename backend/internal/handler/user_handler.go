package handler

import (
	"encoding/json"
	"net/http"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/service"

	"stride-wars-app/ent"

	"github.com/google/uuid"
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Check for 'username' or 'id' query param
	username := r.URL.Query().Get("username")
	idStr := r.URL.Query().Get("id")

	switch {
	case username != "":
		// Fetch by username
		resp, err := h.userService.FindByUsername(r.Context(), username)
		if err != nil {
			h.logger.Error("find user by username failed", zap.Error(err))
			if ent.IsNotFound(err) {
				middleware.WriteError(w, http.StatusNotFound, "user not found")
			} else {
				middleware.WriteError(w, http.StatusBadRequest, err.Error())
			}
			return
		}
		middleware.WriteJSON(w, http.StatusOK, resp)
		return

	case idStr != "":
		// Validate UUID
		userID, err := uuid.Parse(idStr)
		if err != nil {
			middleware.WriteError(w, http.StatusBadRequest, "Invalid UUID format for 'id'")
			return
		}

		// Fetch by ID
		resp, err := h.userService.FindByID(r.Context(), userID)
		if err != nil {
			h.logger.Error("find user by ID failed", zap.Error(err))
			if ent.IsNotFound(err) {
				middleware.WriteError(w, http.StatusNotFound, "user not found")
			} else {
				middleware.WriteError(w, http.StatusBadRequest, err.Error())
			}
			return
		}
		middleware.WriteJSON(w, http.StatusOK, resp)
		return

	default:
		middleware.WriteError(w, http.StatusBadRequest, "Missing 'username' or 'id' query parameter")
	}
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
	if err == nil {
		middleware.WriteJSON(w, http.StatusOK, resp)
		return
	}

	h.logger.Error("update username failed", zap.Error(err))
	switch {
	case ent.IsNotFound(err):
		middleware.WriteError(w, http.StatusNotFound, "user not found")
	default:
		middleware.WriteError(w, http.StatusInternalServerError, "could not update username")
	}

}

// merge get into one func

package handler

import (
	"net/http"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/dto"
	"stride-wars-app/internal/service"
	"stride-wars-app/internal/util"
	"github.com/google/uuid"
	"stride-wars-app/ent"

	"go.uber.org/zap"
)

type ActivityHandler struct {
	activityService *service.ActivityService
	logger          *zap.Logger
}

func NewActivityHandler(activityService *service.ActivityService, logger *zap.Logger) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
		logger:          logger,
	}
}

func (h *ActivityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {

	// Can't use middleware.ParseJSON because it maps int64 into float64 causing inaccurate h3indexes
	// Unmarshal into the request struct
	var req dto.CreateActivityRequest
	if err := util.DecodeJSONBody(r.Body, &req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	activity := dto.CreateActivityRequest{
		UserID:    req.UserID,
		H3Indexes: req.H3Indexes,
		Duration:  req.Duration,
		Distance:  req.Distance,
	}
	resp, err := h.activityService.CreateActivity(r.Context(), activity)
	if err != nil {
		h.logger.Error("create activity failed", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, resp)
}

func (h *ActivityHandler) GetUserActivityStats(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("user_id")

	if idStr == "" {
		middleware.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid UUID format for 'user_id'")
		return
	}

	ActivityStats, err := h.activityService.GetUserActivityStats(r.Context(), userID)
	if err != nil {
		h.logger.Error("find user activities failed", zap.Error(err))
		if ent.IsNotFound(err) {
			middleware.WriteError(w, http.StatusNotFound, "No activities found for this user")
		} else {
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	middleware.WriteJSON(w, http.StatusOK, ActivityStats)
}

package handler

import (
	"net/http"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/service"
	util "stride-wars-app/internal/utils"

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
	var req service.CreateActivityRequest
	if err := util.DecodeJSONBody(r.Body, &req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	activity := service.CreateActivityRequest{
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

package handler

import (
	"encoding/json"
	"net/http"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/service"

	"go.uber.org/zap"
)

type HexLeaderboardHandler struct {
	hexLeaderboardService *service.HexLeaderboardService
	logger                *zap.Logger
}

func NewHexLeaderboardHandler(hexLeaderboardService *service.HexLeaderboardService, logger *zap.Logger) *HexLeaderboardHandler {
	return &HexLeaderboardHandler{
		hexLeaderboardService: hexLeaderboardService,
		logger:                logger,
	}
}

func (h *HexLeaderboardHandler) GetAllLeaderboardsInsideBBBox(w http.ResponseWriter, r *http.Request) {
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
	var req service.GetAllLeaderboardsInsideBBBoxRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	resp, err := h.hexLeaderboardService.GetAllLeaderboardsInsideBBBox(r.Context(), req.BoundingBox)
	if err != nil {
		h.logger.Error("get all leaderboards inside bbox failed", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

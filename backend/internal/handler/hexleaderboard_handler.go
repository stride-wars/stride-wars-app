package handler

import (
	"net/http"
	"strconv"
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

func (h HexLeaderboardHandler) GetAllLeaderboardsInsideBBBox(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	minLatStr := query.Get("min_lat")
	minLngStr := query.Get("min_lng")
	maxLatStr := query.Get("max_lat")
	maxLngStr := query.Get("max_lng")

	// Validate that all required parameters are present
	if minLatStr == "" || minLngStr == "" || maxLatStr == "" || maxLngStr == "" {
		middleware.WriteError(w, http.StatusBadRequest, "Missing required parameters: min_lat, min_lng, max_lat, max_lng")
		return
	}

	// Parse string parameters to float64
	minLat, err := strconv.ParseFloat(minLatStr, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid min_lat parameter")
		return
	}

	minLng, err := strconv.ParseFloat(minLngStr, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid min_lng parameter")
		return
	}

	maxLat, err := strconv.ParseFloat(maxLatStr, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid max_lat parameter")
		return
	}

	maxLng, err := strconv.ParseFloat(maxLngStr, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Invalid max_lng parameter")
		return
	}

	// Create bounding box from parsed parameters
	boundingBox := service.BoundingBox{
		MinLat: minLat,
		MinLng: minLng,
		MaxLat: maxLat,
		MaxLng: maxLng,
	}

	// Call the service with the bounding box
	resp, err := h.hexLeaderboardService.GetAllLeaderboardsInsideBBBox(r.Context(), boundingBox)
	if err != nil {
		h.logger.Error("get all leaderboards inside bbox failed", zap.Error(err))
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

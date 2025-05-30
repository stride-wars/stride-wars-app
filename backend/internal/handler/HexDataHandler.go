package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"stride-wars-app/ent"

	"go.uber.org/zap"
)

type HexDataHandler struct {
	Logger    *zap.Logger
	EntClient *ent.Client
}

func NewHexDataHandler(logger *zap.Logger, entClient *ent.Client) *HexDataHandler {
	return &HexDataHandler{Logger: logger, EntClient: entClient}
}

func (h *HexDataHandler) ReceiveHexData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int64
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Parse hex string ID to int64, base 16
	tempID, err := strconv.ParseInt(req.ID, 16, 64)
	if err != nil {
		http.Error(w, "ID must be a valid hex string", http.StatusBadRequest)
		return
	}

	hex, err := h.EntClient.Hex.Create().
		SetID(tempID).
		Save(r.Context())
	if err != nil {
		if ent.IsConstraintError(err) {
			h.Logger.Info("Hex already exists", zap.Int64("id", req.ID))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"already exists"}`))
			return
		}
		h.Logger.Error("Failed to save hex", zap.Error(err))
		http.Error(w, "Failed to save hex", http.StatusInternalServerError)
		return
	}

	// Convert saved int64 ID back to hex string for logging
	tempID2 := strconv.FormatInt(hex.ID, 16)
	h.Logger.Info("Saved hex", zap.String("id", tempID2))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

package handler

import (
	"encoding/json"
	"net/http"
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

	hex, err := h.EntClient.Hex.Create().
		SetID(req.ID).
		Save(r.Context())
	if err != nil {
		// Check if it's a duplicate key error
		if ent.IsConstraintError(err) {
			h.Logger.Info("Hex already exists", zap.Int64("id", req.ID))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"already exists"}`))
			return
		}
		h.Logger.Error("Failed to save hex", zap.Error(err))
		http.Error(w, "Failed to save hex", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Saved hex", zap.Int64("id", hex.ID))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

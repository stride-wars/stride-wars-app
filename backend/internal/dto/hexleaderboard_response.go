package dto

import (
	"github.com/google/uuid"
)

type TopUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Score  float64   `json:"score"`
}

type HexLeaderboardResponse struct {
	ID       uuid.UUID         `json:"id"`
	H3Index  int64             `json:"h3_index"`
	TopUsers []TopUserResponse `json:"top_users"`
}

type GetAllHexLeaderboardsInsideBBoxResponse struct {
	Leaderboards []HexLeaderboardResponse `json:"leaderboards"`
}
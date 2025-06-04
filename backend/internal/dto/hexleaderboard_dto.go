package dto

import (
	"github.com/google/uuid"
)

type TopUserResponse struct {
	UserID uuid.UUID	`json:"user_id"`
	UserName string 	`json:"user_name"`
	Score  float64   	`json:"score"`
}

type HexLeaderboardResponse struct {
	ID       uuid.UUID         `json:"id"`
	H3Index  string            `json:"h3_index"`
	TopUsers []TopUserResponse `json:"top_users"`
}

type GetAllHexLeaderboardsInsideBBoxResponse struct {
	Leaderboards []HexLeaderboardResponse `json:"leaderboards"`
}

type GlobalLeaderboardEntry struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	TopCount int       `json:"top_count"`
}

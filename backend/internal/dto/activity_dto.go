package dto

import (
	"github.com/google/uuid"
)

type CreateActivityRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	Duration  float64   `json:"duration"` // in seconds
	Distance  float64   `json:"distance"` // in meters
	H3Indexes []int64   `json:"h3_indexes"`
}

type CreateActivityResponse struct {
	ID        uuid.UUID `json:"activity_id"`
	UserID    uuid.UUID `json:"user_id"`
	Duration  float64   `json:"duration"` // in seconds
	Distance  float64   `json:"distance"` // in meters
	H3Indexes []int64   `json:"h3_indexes"`
}

type GetUserActivityStatsResponse struct {
	HexesVisited int64	`json:"hexes_visited"`
	ActivitiesRecorded int64	`json:"activities_recorded"`
	DistanceCovered float64 `json:"distance_covered"` // in meters
	WeeklyActivities []int64 `json:"weekly_activities"`
}

// internal/mappers/hexleaderboard_mapper.go
package mappers

import (
	"sort"

	"stride-wars-app/ent"
	"stride-wars-app/internal/dto"
)

func MapHexLeaderboardsToResponse(hexLeaderboards []*ent.HexLeaderboard) *dto.GetAllHexLeaderboardsInsideBBoxResponse {
	leaderboards := make([]dto.HexLeaderboardResponse, 0, len(hexLeaderboards))

	for _, hexLeaderboard := range hexLeaderboards {
		topUsers := make([]dto.TopUserResponse, 0, len(hexLeaderboard.TopUsers))
		for _, user := range hexLeaderboard.TopUsers {
			topUsers = append(topUsers, dto.TopUserResponse{
				UserID: user.UserID,
				Score:  user.Score,
			})
		}

		leaderboards = append(leaderboards, dto.HexLeaderboardResponse{
			ID:       hexLeaderboard.ID,
			H3Index:  hexLeaderboard.H3Index,
			TopUsers: topUsers,
		})
	}

	sort.Slice(leaderboards, func(i, j int) bool {
		return leaderboards[i].H3Index < leaderboards[j].H3Index
	})

	return &dto.GetAllHexLeaderboardsInsideBBoxResponse{
		Leaderboards: leaderboards,
	}
}

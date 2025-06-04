package repository

import (
	"context"
	"sort"
	"stride-wars-app/ent"
	entHexLeaderboard "stride-wars-app/ent/hexleaderboard"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/dto"

	"github.com/google/uuid"
)

type HexLeaderboardRepository struct {
	client *ent.Client
}

func NewHexLeaderboardRepository(client *ent.Client) HexLeaderboardRepository {
	return HexLeaderboardRepository{client: client}
}
func (r HexLeaderboardRepository) FindByID(ctx context.Context, id uuid.UUID) (*ent.HexLeaderboard, error) {
	return r.client.HexLeaderboard.Query().Where(entHexLeaderboard.IDEQ(id)).First(ctx)
}
func (r HexLeaderboardRepository) FindByH3Index(ctx context.Context, hexID int64) (*ent.HexLeaderboard, error) {
	return r.client.HexLeaderboard.Query().Where(entHexLeaderboard.H3IndexEQ(hexID)).First(ctx)
}
func (r HexLeaderboardRepository) CreateHexLeaderboard(ctx context.Context, hexLeaderboard *model.HexLeaderboard) (*ent.HexLeaderboard, error) {
	return r.client.HexLeaderboard.Create().
		SetH3Index(hexLeaderboard.H3Index).
		SetTopUsers(hexLeaderboard.TopUsers).
		Save(ctx)
}
func (r HexLeaderboardRepository) UpdateHexLeaderboard(ctx context.Context, hexLeaderboard *model.HexLeaderboard) (int, error) {
	return r.client.HexLeaderboard.Update().Where(entHexLeaderboard.IDEQ(hexLeaderboard.ID)).SetTopUsers(hexLeaderboard.TopUsers).Save(ctx)
}
func (r HexLeaderboardRepository) FindByH3Indexes(ctx context.Context, h3Indexes []int64) ([]*ent.HexLeaderboard, error) {
	return r.client.HexLeaderboard.Query().Where(entHexLeaderboard.H3IndexIn(h3Indexes...)).All(ctx)
}

// Return users position in a particular hex's leaderboard, returns nil if the user is not in the leaderboard / in case of an error
func (r HexLeaderboardRepository) GetUserPositionInLeaderboard(ctx context.Context, hexID int64, userID uuid.UUID) (*int, error) {
	hexLeaderboard, err := r.FindByH3Index(ctx, hexID)
	if err != nil {
		return nil, err
	}
	if hexLeaderboard == nil {
		return nil, nil
	}

	topUsers := hexLeaderboard.TopUsers
	for idx, user := range topUsers {
		if user.UserID == userID {
			pos := idx + 1
			return &pos, nil
		}
	}
	return nil, nil
}

func (r HexLeaderboardRepository) GetGlobalHexLeaderboard(ctx context.Context) ([]dto.GlobalLeaderboardEntry, error) {
	leaderboards, err := r.client.HexLeaderboard.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	userCounts := make(map[uuid.UUID]int)
	for _, lb := range leaderboards {
		if len(lb.TopUsers) > 0 {
			userCounts[lb.TopUsers[0].UserID]++
		}
	}

	var entries []dto.GlobalLeaderboardEntry
	for userID, count := range userCounts {
		entries = append(entries, dto.GlobalLeaderboardEntry{UserID: userID, TopCount: count})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].TopCount > entries[j].TopCount
	})

	if len(entries) > 10 {
		entries = entries[:10]
	}
	return entries, nil
}

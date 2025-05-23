package repository

import (
	"context"
	"stride-wars-app/ent"
	entHexLeaderboard "stride-wars-app/ent/hexleaderboard"
	"stride-wars-app/ent/model"

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

// Return users position in a particualr hex's leaderboard, returns -1 if the user is not in the leaderboard / in case of an error
func (r HexLeaderboardRepository) GetUserPositionInLeaderboard(ctx context.Context, hexID int64, userID uuid.UUID) (int, error) {
	hexLeaderboard, err := r.FindByH3Index(ctx, hexID)
	if err != nil {
		return -1, err
	}
	if hexLeaderboard == nil {
		return -1, nil
	}

	topUsers := hexLeaderboard.TopUsers
	for i, user := range topUsers {
		if user.UserID == userID {
			return i + 1, nil
		}
	}
	return -1, nil
}

package service	
import (
	"context"	
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sort"
)
type HexLeaderboardService struct {
	hexLeaderboardRepository repository.HexLeaderboardRepository
	hexInfluenceRepository repository.HexInfluenceRepository
	logger     *zap.Logger
}
func NewHexLeaderboardService(hexLeaderboardRepository repository.HexLeaderboardRepository, hexInfluenceRepository repository.HexInfluenceRepository, logger *zap.Logger) *HexLeaderboardService {
	return &HexLeaderboardService{
		hexLeaderboardRepository: hexLeaderboardRepository,
		hexInfluenceRepository:   hexInfluenceRepository,
		logger:     logger,
	}
}
func (hls *HexLeaderboardService) FindByID(ctx context.Context, id uuid.UUID) (*ent.HexLeaderboard, error) {
	return hls.hexLeaderboardRepository.FindByID(ctx, id)
}
//find by h3 index
func (hls *HexLeaderboardService) FindByH3Index(ctx context.Context, h3_index int64) (*ent.HexLeaderboard, error) {
	return hls.hexLeaderboardRepository.FindByH3Index(ctx, h3_index)
}
func (hls *HexLeaderboardService) CreateHexLeaderboard(ctx context.Context, hexLeaderboard *model.HexLeaderboard) (*ent.HexLeaderboard, error) {
	return hls.hexLeaderboardRepository.CreateHexLeaderboard(ctx, hexLeaderboard)
}
func (hls *HexLeaderboardService) UpdateHexLeaderboard(ctx context.Context, hexLeaderboard *model.HexLeaderboard) (int, error) {
	return hls.hexLeaderboardRepository.UpdateHexLeaderboard(ctx, hexLeaderboard)
}
func (hls *HexLeaderboardService) FindByH3Indexes(ctx context.Context, h3Indexes []int64) ([]*ent.HexLeaderboard, error) {
	return hls.hexLeaderboardRepository.FindByH3Indexes(ctx, h3Indexes) 
}
// Ads a given user to the leaderboard of a hexagon with the given hexID - if the user has enough points to go into top 5.
func (hls *HexLeaderboardService) AddUserToLeaderboard(ctx context.Context, hexID int64, userID uuid.UUID) (int, error) {
	hexLeaderboard, err := hls.hexLeaderboardRepository.FindByH3Index(ctx, hexID)
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, nil 
		}
		return 0, err
	}

	hexInfluence, err := hls.hexInfluenceRepository.FindByUserIDAndHexID(ctx, userID, hexID)
	if err != nil {
		return 0, err
	}
	if hexInfluence == nil {
		return 0, nil
	}
	score := hexInfluence.Score

	newTopUsers := make([]model.TopUser, 0, len(hexLeaderboard.TopUsers))
	addedOrUpdated := false

	for _, u := range hexLeaderboard.TopUsers {
		if u.UserID == userID {
			if u.Score != score {
				newTopUsers = append(newTopUsers, model.TopUser{UserID: userID, Score: score})
				addedOrUpdated = true
			} else {
				newTopUsers = append(newTopUsers, u)
			}
		} else {
			newTopUsers = append(newTopUsers, u)
		}
	}

	if !addedOrUpdated {
		newTopUsers = append(newTopUsers, model.TopUser{UserID: userID, Score: score})
	}

	sort.Slice(newTopUsers, func(i, j int) bool {
		return newTopUsers[i].Score > newTopUsers[j].Score
	})

	if len(newTopUsers) > 5 {
		newTopUsers = newTopUsers[:5]
	}

	inTop := false
	for _, u := range newTopUsers {
		if u.UserID == userID {
			inTop = true
			break
		}
	}
	if !inTop {
		return 0, nil 
	}

	hexLeaderboard.TopUsers = newTopUsers
	updatedModel := &model.HexLeaderboard{
		ID:        hexLeaderboard.ID,
		H3Index:   hexLeaderboard.H3Index,
		TopUsers:  newTopUsers,
	}

	_, err = hls.hexLeaderboardRepository.UpdateHexLeaderboard(ctx, updatedModel)
	if err != nil {
		return 0, err
	}

	for i, u := range newTopUsers {
		if u.UserID == userID {
			return i + 1, nil
		}
	}

	return 0, nil 
}


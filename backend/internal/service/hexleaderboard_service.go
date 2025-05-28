package service

import (
	"context"
	"sort"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"

	"github.com/google/uuid"
	"github.com/uber/h3-go/v4"
	"go.uber.org/zap"
)

type BoundingBox struct {
	MinLat float64 `json:"min_lat"`
	MinLng float64 `json:"min_lng"`
	MaxLat float64 `json:"max_lat"`
	MaxLng float64 `json:"max_lng"`
}

type GetAllLeaderboardsInsideBBBoxRequest struct {
	BoundingBox BoundingBox `json:"bounding_box"`
}
type GetAllLeaderboardsInsideBBBoxResponse struct {
	Leaderboards []*ent.HexLeaderboard `json:"leaderboards"`
}

type HexLeaderboardService struct {
	hexLeaderboardRepository repository.HexLeaderboardRepository
	hexInfluenceRepository   repository.HexInfluenceRepository
	logger                   *zap.Logger
}

func NewHexLeaderboardService(hexLeaderboardRepository repository.HexLeaderboardRepository, hexInfluenceRepository repository.HexInfluenceRepository, logger *zap.Logger) *HexLeaderboardService {
	return &HexLeaderboardService{
		hexLeaderboardRepository: hexLeaderboardRepository,
		hexInfluenceRepository:   hexInfluenceRepository,
		logger:                   logger,
	}
}
func (hls *HexLeaderboardService) FindByID(ctx context.Context, id uuid.UUID) (*ent.HexLeaderboard, error) {
	return hls.hexLeaderboardRepository.FindByID(ctx, id)
}

// find by h3 index
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
// Return users position in the leaderboard - nil otherwise
func (hls *HexLeaderboardService) AddUserToLeaderboard(ctx context.Context, hexID int64, userID uuid.UUID) (*int, error) {
	hexLeaderboard, err := hls.hexLeaderboardRepository.FindByH3Index(ctx, hexID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	hexInfluence, err := hls.hexInfluenceRepository.FindByUserIDAndHexID(ctx, userID, hexID)
	if err != nil {
		return nil, err
	}
	if hexInfluence == nil {
		return nil, nil
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
		return nil, nil
	}

	hexLeaderboard.TopUsers = newTopUsers
	updatedModel := &model.HexLeaderboard{
		ID:       hexLeaderboard.ID,
		H3Index:  hexLeaderboard.H3Index,
		TopUsers: newTopUsers,
	}

	_, err = hls.hexLeaderboardRepository.UpdateHexLeaderboard(ctx, updatedModel)
	if err != nil {
		return nil, err
	}

	for idx, u := range newTopUsers {
		if u.UserID == userID {
			pos := idx + 1
			return &pos, nil
		}
	}

	return nil, nil
}

// func AddUserToLeaderboardOrCreate
func (hls *HexLeaderboardService) AddUserToLeaderboardOrCreateLeaderboard(ctx context.Context, hexID int64, userID uuid.UUID) (*int, error) {
	// Try to find the leaderboard
	_, err := hls.hexLeaderboardRepository.FindByH3Index(ctx, hexID)
	if err != nil {
		if ent.IsNotFound(err) {
			// Attempt to get influence score first
			hexInfluence, infErr := hls.hexInfluenceRepository.FindByUserIDAndHexID(ctx, userID, hexID)
			if infErr != nil {
				return nil, infErr
			}
			if hexInfluence == nil {
				return nil, nil
			}

			// Create leaderboard with current user as first entry
			leaderboard := &model.HexLeaderboard{
				H3Index: hexID,
				TopUsers: []model.TopUser{
					{UserID: userID, Score: hexInfluence.Score},
				},
			}

			_, err = hls.hexLeaderboardRepository.CreateHexLeaderboard(ctx, leaderboard)
			if err != nil {
				return nil, err
			}
			pos := 1

			return &pos, nil // user is the first and only one in the leaderboard
		}
		return nil, err // real error
	}

	// If leaderboard exists, attempt to add user
	position, err := hls.AddUserToLeaderboard(ctx, hexID, userID)
	if err != nil {
		return nil, err
	}
	return position, nil
}

// Return users position in a particular hex's leaderboard, returns nil if the user is not in the leaderboard / in case of an error
func (hls *HexLeaderboardService) GetUserPositionInLeaderboard(ctx context.Context, hexID int64, userID uuid.UUID) (*int, error) {
	return hls.hexLeaderboardRepository.GetUserPositionInLeaderboard(ctx, hexID, userID)
}

// returns all existing hex leaderboards inside a given bounding box
func (hls *HexLeaderboardService) GetAllLeaderboardsInsideBBBox(ctx context.Context, bbox BoundingBox) ([]*ent.HexLeaderboard, error) {

	verts := h3.GeoLoop{
		{Lat: bbox.MinLat, Lng: bbox.MinLng},
		{Lat: bbox.MinLat, Lng: bbox.MaxLng},
		{Lat: bbox.MaxLat, Lng: bbox.MaxLng},
		{Lat: bbox.MaxLat, Lng: bbox.MinLng},
	}

	// Build a GeoPolygon with no holes.
	poly := h3.GeoPolygon{
		GeoLoop: verts,
		Holes:   nil,
	}

	h3Cells, err := h3.PolygonToCells(poly, 9) // 9 is the resolution, adjust as needed
	if err != nil {
		hls.logger.Error("Failed to convert polygon to H3 cells", zap.Error(err))
		return nil, err
	}
	h3Indexes := make([]int64, len(h3Cells))
	for i, cell := range h3Cells {
		h3Indexes[i] = int64(cell)
	}
	// Fetch all hex leaderboards for the given H3 indexes
	hexLeaderboards, err := hls.hexLeaderboardRepository.FindByH3Indexes(ctx, h3Indexes)
	if err != nil {
		hls.logger.Error("Failed to fetch hex leaderboards by H3 indexes", zap.Error(err))
		return nil, err
	}
	return hexLeaderboards, nil
}

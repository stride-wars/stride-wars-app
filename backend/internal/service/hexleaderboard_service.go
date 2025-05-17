package service	
import (
	"context"	
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)
type HexLeaderboardService struct {
	repository repository.HexLeaderboardRepository
	logger     *zap.Logger
}
func NewHexLeaderboardService(repository repository.HexLeaderboardRepository, logger *zap.Logger) *HexLeaderboardService {
	return &HexLeaderboardService{
		repository: repository,
		logger:     logger,
	}
}
func (hls *HexLeaderboardService) FindByID(ctx context.Context, id uuid.UUID) (*ent.HexLeaderboard, error) {
	return hls.repository.FindByID(ctx, id)
}
//find by h3 index
func (hls *HexLeaderboardService) FindByH3Index(ctx context.Context, h3_index int64) ([]*ent.HexLeaderboard, error) {
	return hls.repository.FindByH3Index(ctx, h3_index)
}
func (hls *HexLeaderboardService) CreateHexLeaderboard(ctx context.Context, hexLeaderboard *model.HexLeaderboard) (*ent.HexLeaderboard, error) {
	return hls.repository.CreateHexLeaderboard(ctx, hexLeaderboard)
}
func (hls *HexLeaderboardService) UpdateHexLeaderboard(ctx context.Context, hexLeaderboard *model.HexLeaderboard) (int, error) {
	return hls.repository.UpdateHexLeaderboard(ctx, hexLeaderboard)
}
func (hls *HexLeaderboardService) FindByH3Indexes(ctx context.Context, h3Indexes []int64) ([]*ent.HexLeaderboard, error) {
	return hls.repository.FindByH3Indexes(ctx, h3Indexes)
}
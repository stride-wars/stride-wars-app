package service

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/internal/repository"

	"go.uber.org/zap"
)

type HexService struct {
	repository repository.HexRepository
	logger     *zap.Logger
}

func NewHexService(repository repository.HexRepository, logger *zap.Logger) *HexService {
	return &HexService{
		repository: repository,
		logger:     logger,
	}
}

func (hs *HexService) FindByID(ctx context.Context, h3_index int64) (*ent.Hex, error) {
	return hs.repository.FindByID(ctx, h3_index)
}
func (hs *HexService) FindByIDs(ctx context.Context, h3_indexes []int64) ([]*ent.Hex, error) {
	return hs.repository.FindByIDs(ctx, h3_indexes)
}
func (hs *HexService) CreateHex(ctx context.Context, h3_index int64) (*ent.Hex, error) {
	return hs.repository.CreateHex(ctx, h3_index)
}
func (hs *HexService) CreateHexes(ctx context.Context, hexes []*ent.Hex) ([]*ent.Hex, error) {
	return hs.repository.CreateHexes(ctx, hexes)
}

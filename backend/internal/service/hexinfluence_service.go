package service

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HexInfluenceService struct {
	repository repository.HexInfluenceRepository
	logger     *zap.Logger
}

func NewHexInfluenceService(repository repository.HexInfluenceRepository, logger *zap.Logger) *HexInfluenceService {
	return &HexInfluenceService{
		repository: repository,
		logger:     logger,
	}
}
func (his *HexInfluenceService) FindByID(ctx context.Context, id uuid.UUID) (*ent.HexInfluence, error) {
	return his.repository.FindByID(ctx, id)
}
func (his *HexInfluenceService) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.HexInfluence, error) {
	return his.repository.FindByIDs(ctx, ids)
}
func (his *HexInfluenceService) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.HexInfluence, error) {
	return his.repository.FindByUserID(ctx, userID)
}
func (his *HexInfluenceService) FindByHexID(ctx context.Context, h3_index int64) ([]*ent.HexInfluence, error) {
	return his.repository.FindByHexID(ctx, h3_index)
}
func (his *HexInfluenceService) CreateHexInfluence(ctx context.Context, hexInfluence *model.HexInfluence) (*ent.HexInfluence, error) {
	return his.repository.CreateHexInfluence(ctx, hexInfluence)
}
func (his *HexInfluenceService) UpdateHexInfluence(ctx context.Context, userID uuid.UUID, hexID int64) (int, error) {
	return his.repository.UpdateHexInfluence(ctx, userID, hexID)
}
func (his *HexInfluenceService) FindByUserIDAndHexID(ctx context.Context, userID uuid.UUID, hexID int64) (*ent.HexInfluence, error) {
	return his.repository.FindByUserIDAndHexID(ctx, userID, hexID)
}
func (his *HexInfluenceService) UpdateHexInfluences(ctx context.Context, userID uuid.UUID, hexIDs []int64) (int, error) {
	return his.repository.UpdateHexInfluences(ctx, userID, hexIDs)
}
func (his *HexInfluenceService) UpdateOrCreateHexInfluence(ctx context.Context, userID uuid.UUID, hexID int64) (*ent.HexInfluence, error) {
	return his.repository.UpdateOrCreateHexInfluence(ctx, userID, hexID)
}
func (his *HexInfluenceService) UpdateOrCreateHexInfluences(ctx context.Context, userID uuid.UUID, hexIDs []int64) ([]*ent.HexInfluence, error) {
	return his.repository.UpdateOrCreateHexInfluences(ctx, userID, hexIDs)
}

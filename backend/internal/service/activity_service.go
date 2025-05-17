package service

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ActivityService struct {
	repository             repository.ActivityRepository
	hexInfluenceRepository repository.HexInfluenceRepository
	logger                 *zap.Logger
}
func NewActivityService(repository repository.ActivityRepository, hexInfluenceRepository repository.HexInfluenceRepository, logger *zap.Logger) *ActivityService {
	return &ActivityService{
		repository:             repository,
		hexInfluenceRepository: hexInfluenceRepository,
		logger:                 logger,
	}
}


func (as *ActivityService) FindByID(ctx context.Context, uuid uuid.UUID) (*ent.Activity, error) {
	return as.repository.FindByID(ctx, uuid)
}
func (as *ActivityService) FindByIDs(ctx context.Context, uuids []uuid.UUID) ([]*ent.Activity, error) {
	return as.repository.FindByIDs(ctx, uuids)
}
func (as *ActivityService) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.Activity, error) {
	return as.repository.FindByUserID(ctx, userID)
}
func (as *ActivityService) CreateActivity(ctx context.Context, activity *model.Activity) (*ent.Activity, error) {

	created_activity, err := as.repository.CreateActivity(ctx, activity)
	if err != nil {
		as.logger.Error("failed to create activity", zap.Error(err))
		return nil, err
	}
	his := NewHexInfluenceService(as.hexInfluenceRepository, as.logger) 
	// TO DO: finish creating the whole chain

}

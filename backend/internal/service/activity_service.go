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
	hexRepository            repository.HexRepository 
	hexInfluenceRepository repository.HexInfluenceRepository
	hexLeaderboardRepository repository.HexLeaderboardRepository
	logger                 *zap.Logger
}

func NewActivityService(
	activityRepo repository.ActivityRepository,
	hexRepo repository.HexRepository, 
	hexInfluenceRepo repository.HexInfluenceRepository,
	hexLeaderboardRepo repository.HexLeaderboardRepository, 
	logger *zap.Logger,
) *ActivityService {
	return &ActivityService{
		repository:               activityRepo,
		hexRepository:            hexRepo, 
		hexInfluenceRepository:   hexInfluenceRepo,
		hexLeaderboardRepository: hexLeaderboardRepo, 
		logger:                   logger,
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



func (as *ActivityService) CreateActivity(ctx context.Context, activityInput *model.Activity) (*ent.Activity, error) {

	createdActivity, err := as.repository.CreateActivity(ctx, activityInput)
	if err != nil {
		return nil, err
	}

	userID := activityInput.UserID
	h3Indexes := activityInput.H3Indexes 

	if len(h3Indexes) == 0 {
		as.logger.Info("Activity contains no H3 indexes; skipping hex influence and leaderboard updates.", zap.Stringer("activityID", createdActivity.ID))
		return createdActivity, nil
	}

	hexService := NewHexService(as.hexRepository, as.logger)
	hexInfluenceService := NewHexInfluenceService(as.hexInfluenceRepository, as.logger)
	hexLeaderboardService := NewHexLeaderboardService(as.hexLeaderboardRepository,as.hexInfluenceRepository, as.logger)

	as.logger.Debug("Checking and creating hexes if necessary", zap.Any("h3Indexes", h3Indexes))
	existingHexEntities, err := hexService.FindByIDs(ctx, h3Indexes) 
	if err != nil {
		as.logger.Warn("Failed to pre-fetch existing hexes by IDs. Will attempt creation individually.",
			zap.Error(err), zap.Any("h3Indexes", h3Indexes))
	}

	existingHexMap := make(map[int64]bool)
	for _, eh := range existingHexEntities {
		existingHexMap[eh.ID] = true
	}

	for _, h3Index := range h3Indexes {
		if _, exists := existingHexMap[h3Index]; !exists {
			as.logger.Info("Hex not found in database, attempting to create.", zap.Int64("h3Index", h3Index))
			_, createHexErr := hexService.CreateHex(ctx, h3Index)
			if createHexErr != nil {
				as.logger.Error("Failed to create hex, or it was created concurrently by another process.",
					zap.Error(createHexErr), zap.Int64("h3Index", h3Index))
			} else {
				as.logger.Info("Successfully created new hex in database.", zap.Int64("h3Index", h3Index))
			}
		}
	}
	as.logger.Debug("Hex existence check and creation phase complete.")

	for _, h3Index := range h3Indexes {

		// Update or create hex influence
		_, err := hexInfluenceService.UpdateOrCreateHexInfluence(ctx, userID, h3Index)
		if err != nil {
			as.logger.Error("Failed to update or create hex influence.", zap.Error(err), zap.Int64("h3Index", h3Index))
			continue
		}

		// Add user to leaderboard
		_, err = hexLeaderboardService.AddUserToLeaderboard(ctx, h3Index, userID)
		if err != nil {
			as.logger.Error("Failed to add user to leaderboard.", zap.Error(err), zap.Int64("h3Index", h3Index))
			continue
		}
		as.logger.Info("Successfully added user to leaderboard.", zap.Int64("h3Index", h3Index))
	}
	as.logger.Info("Finished processing all H3 indexes for activity.", zap.Stringer("activityID", createdActivity.ID))

	return createdActivity, nil
}

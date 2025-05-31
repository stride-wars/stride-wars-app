package service

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/constants"
	"stride-wars-app/internal/repository"

	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/uber/h3-go/v4"
	"go.uber.org/zap"
)

type CreateActivityRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	Duration  float64   `json:"duration"` // in seconds
	Distance  float64   `json:"distance"` // in meters
	H3Indexes []int64   `json:"h3_indexes"`
}

type CreateActivityResponse struct {
	ID        uuid.UUID `json:"activity_id"`
	UserID    uuid.UUID `json:"user_id"`
	Duration  float64   `json:"duration"` // in seconds
	Distance  float64   `json:"distance"` // in meters
	H3Indexes []int64   `json:"h3_indexes"`
}

type ActivityService struct {
	repository               repository.ActivityRepository
	hexRepository            repository.HexRepository
	hexInfluenceRepository   repository.HexInfluenceRepository
	hexLeaderboardRepository repository.HexLeaderboardRepository
	userRepository           repository.UserRepository
	userService              UserService
	logger                   *zap.Logger
}

func NewActivityService(
	activityRepo repository.ActivityRepository,
	hexRepo repository.HexRepository,
	hexInfluenceRepo repository.HexInfluenceRepository,
	hexLeaderboardRepo repository.HexLeaderboardRepository,
	userRepo repository.UserRepository,
	userService UserService,
	logger *zap.Logger,
) *ActivityService {
	return &ActivityService{
		repository:               activityRepo,
		hexRepository:            hexRepo,
		hexInfluenceRepository:   hexInfluenceRepo,
		hexLeaderboardRepository: hexLeaderboardRepo,
		userRepository:           userRepo,
		userService:              userService,
		logger:                   logger,
	}
}

func (a *ActivityService) validateCreateActivity(req CreateActivityRequest) error {
	if (req.UserID == uuid.Nil || req.UserID == uuid.UUID{}) {
		return errors.New("UserID is required")
	}

	if req.Duration <= 0 {
		return errors.New("duration must be positive")
	}

	if req.Distance <= 0 {
		return errors.New("distance must be positive")
	}

	if len(req.H3Indexes) == 0 {
		return errors.New("at least one H3 index is required")
	}
	for _, h3Index := range req.H3Indexes {
		cell := h3.Cell(h3Index)
		if !cell.IsValid() {
			return errors.New("invalid H3 index: " + strconv.FormatInt(h3Index, 10))
		}

		if cell.Resolution() != constants.DefaultHexResolution {
			return errors.New("H3 index " + strconv.FormatInt(h3Index, 10) + " is not at resolution " + strconv.Itoa(constants.DefaultHexResolution))
		}
	}

	return nil
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

func (as *ActivityService) CreateActivity(ctx context.Context, req CreateActivityRequest) (*CreateActivityResponse, error) {
	if err := as.validateCreateActivity(req); err != nil {
		return nil, err
	}

	activityInput := &model.Activity{
		UserID:    req.UserID,
		Duration:  req.Duration,
		Distance:  req.Distance,
		H3Indexes: req.H3Indexes,
	}
	if len(activityInput.H3Indexes) == 0 {
		return nil, errors.New("activity must contain at least one H3 index")
	}
	// validate if user exists
	_, err := as.userService.FindByID(ctx, activityInput.UserID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, err
	}

	createdActivity, err := as.repository.CreateActivity(ctx, activityInput)
	if err != nil {
		return nil, err
	}

	userID := activityInput.UserID
	h3Indexes := activityInput.H3Indexes

	hexService := NewHexService(as.hexRepository, as.logger)
	hexInfluenceService := NewHexInfluenceService(as.hexInfluenceRepository, as.logger)
	hexLeaderboardService := NewHexLeaderboardService(as.hexLeaderboardRepository, as.hexInfluenceRepository, as.logger)

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
		_, err = hexLeaderboardService.AddUserToLeaderboardOrCreateLeaderboard(ctx, h3Index, userID)
		if err != nil {
			as.logger.Error("Failed to add user to leaderboard.", zap.Error(err), zap.Int64("h3Index", h3Index))
			continue
		}
		as.logger.Info("Successfully added user to leaderboard.", zap.Int64("h3Index", h3Index))
	}
	as.logger.Info("Finished processing all H3 indexes for activity.", zap.Stringer("activityID", createdActivity.ID))

	return &CreateActivityResponse{
		ID:        createdActivity.ID,
		UserID:    createdActivity.UserID,
		Duration:  createdActivity.DurationSeconds,
		Distance:  createdActivity.DistanceMeters,
		H3Indexes: createdActivity.H3Indexes,
	}, nil
}

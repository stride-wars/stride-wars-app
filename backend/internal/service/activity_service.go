package service

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/dto"
	"stride-wars-app/internal/hex/hexconsts"
	"stride-wars-app/internal/repository"

	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/uber/h3-go/v4"
	"go.uber.org/zap"
)

type ActivityService struct {
	repository            repository.ActivityRepository
	HexService            *HexService
	HexInfluenceService   *HexInfluenceService
	HexLeaderboardService *HexLeaderboardService
	UserService           *UserService
	logger                *zap.Logger
}

func NewActivityService(
	activityRepo repository.ActivityRepository,
	hexInfluenceRepo repository.HexInfluenceRepository,
	hexLeaderboardRepo repository.HexLeaderboardRepository,
	hexRepo repository.HexRepository,
	userService *UserService, // Fixed: pass already constructed service
	logger *zap.Logger,
) *ActivityService {
	return &ActivityService{
		repository:            activityRepo,
		HexService:            NewHexService(hexRepo, logger),
		HexInfluenceService:   NewHexInfluenceService(hexInfluenceRepo, logger),
		HexLeaderboardService: NewHexLeaderboardService(hexLeaderboardRepo, hexInfluenceRepo, logger),
		UserService:           userService, // Fixed: use passed-in service
		logger:                logger,
	}
}

func (a *ActivityService) validateCreateActivity(req dto.CreateActivityRequest) error {
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
		cell := h3.Cell(h3.IndexFromString(h3Index))
		if !cell.IsValid() {
			return errors.New("invalid H3 index: " + h3Index)
		}

		if cell.Resolution() != hexconsts.DefaultHexResolution {
			return errors.New("H3 index " + h3Index + " is not at resolution " + strconv.Itoa(hexconsts.DefaultHexResolution))
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

func (as *ActivityService) CreateActivity(ctx context.Context, req dto.CreateActivityRequest) (*dto.CreateActivityResponse, error) {
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
	_, err := as.UserService.FindByID(ctx, activityInput.UserID)
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
	userName := "Jan" // TODO: change this to a valid value

	as.logger.Debug("Checking and creating hexes if necessary", zap.Any("h3Indexes", h3Indexes))
	existingHexEntities, err := as.HexService.FindByIDs(ctx, h3Indexes)
	if err != nil {
		as.logger.Warn("Failed to pre-fetch existing hexes by IDs. Will attempt creation individually.",
			zap.Error(err), zap.Any("h3Indexes", h3Indexes))
	}

	existingHexMap := make(map[string]bool)
	for _, eh := range existingHexEntities {
		existingHexMap[eh.ID] = true
	}

	for _, h3Index := range h3Indexes {
		if _, exists := existingHexMap[h3Index]; !exists {
			as.logger.Info("Hex not found in database, attempting to create.", zap.String("h3Index", h3Index))
			_, createHexErr := as.HexService.CreateHex(ctx, h3Index)
			if createHexErr != nil {
				as.logger.Error("Failed to create hex, or it was created concurrently by another process.",
					zap.Error(createHexErr), zap.String("h3Index", h3Index))
			} else {
				as.logger.Info("Successfully created new hex in database.", zap.String("h3Index", h3Index))
			}
		}
	}
	as.logger.Debug("Hex existence check and creation phase complete.")

	for _, h3Index := range h3Indexes {

		// Update or create hex influence
		_, err := as.HexInfluenceService.UpdateOrCreateHexInfluence(ctx, userID, h3Index)
		if err != nil {
			as.logger.Error("Failed to update or create hex influence.", zap.Error(err), zap.String("h3Index", h3Index))
			continue
		}

		// Add user to leaderboard
		_, err = as.HexLeaderboardService.AddUserToLeaderboardOrCreateLeaderboard(ctx, h3Index, userID, userName)
		if err != nil {
			as.logger.Error("Failed to add user to leaderboard.", zap.Error(err), zap.String("h3Index", h3Index))
			continue
		}
		as.logger.Info("Successfully added user to leaderboard.", zap.String("h3Index", h3Index))
	}
	as.logger.Info("Finished processing all H3 indexes for activity.", zap.Stringer("activityID", createdActivity.ID))

	return &dto.CreateActivityResponse{
		ID:        createdActivity.ID,
		UserID:    createdActivity.UserID,
		Duration:  createdActivity.DurationSeconds,
		Distance:  createdActivity.DistanceMeters,
		H3Indexes: createdActivity.H3Indexes,
	}, nil
}

func (as *ActivityService) GetUserActivityStats(ctx context.Context, userID uuid.UUID) (*dto.GetUserActivityStatsResponse, error) {
	activities, err := as.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return &dto.GetUserActivityStatsResponse{
			HexesVisited:       0,
			ActivitiesRecorded: 0,
			DistanceCovered:    0,
			WeeklyActivities:   make([]int64, 7),
		}, nil
	}

	stats := &dto.GetUserActivityStatsResponse{
		HexesVisited:       0,
		ActivitiesRecorded: int64(len(activities)),
		DistanceCovered:    0,
		WeeklyActivities:   make([]int64, 7),
	}

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for _, activity := range activities {
		stats.DistanceCovered += activity.DistanceMeters
		stats.HexesVisited += int64(len(activity.H3Indexes))

		// Calculate how many days ago this activity was
		daysAgo := int(startOfToday.Sub(activity.CreatedAt).Hours() / 24)
		if daysAgo >= 0 && daysAgo < 7 {
			stats.WeeklyActivities[6-daysAgo]++ // Reverse so index 6 is today
		}
	}

	return stats, nil
}

package service

import (
	"stride-wars-app/internal/repository"

	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

type Services struct {
	UserService           *UserService
	AuthService           *AuthService
	ActivityService       *ActivityService
	HexService            *HexService
	HexLeaderboardService *HexLeaderboardService
	HexInfluenceService   *HexInfluenceService
}

func Provide(repositories *repository.Repositories, supabaseClient *supabase.Client, logger *zap.Logger) *Services {
	userService := NewUserService(repositories.UserRepository, logger)

	return &Services{
		UserService: userService,
		AuthService: NewAuthService(supabaseClient, logger, userService),
		ActivityService: NewActivityService(repositories.ActivityRepository,
			repositories.HexRepository,
			repositories.HexInfluenceRepository,
			repositories.HexLeaderboardRepository,
			repositories.UserRepository,
			*userService,
			logger),
		HexService: NewHexService(repositories.HexRepository, logger),
		HexLeaderboardService: NewHexLeaderboardService(repositories.HexLeaderboardRepository,
			repositories.HexInfluenceRepository,
			logger),
		HexInfluenceService: NewHexInfluenceService(repositories.HexInfluenceRepository, logger),
	}
}

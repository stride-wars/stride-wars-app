package service

import (
	"stride-wars-app/ent"
	"stride-wars-app/internal/repository"

	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

type Service struct {
	UserService *UserService
	AuthService *AuthService
	ActivityService *ActivityService
	HexService *HexService
	HexLeaderboardService *HexLeaderboardService
	HexInfluenceService *HexInfluenceService
}

func Provide(entClient *ent.Client, supabaseClient *supabase.Client, logger *zap.Logger) *Service {
	return &Service{
		UserService: NewUserService(repository.NewUserRepository(entClient), logger),
		AuthService: NewAuthService(supabaseClient, logger),
		ActivityService: NewActivityService(repository.NewActivityRepository(entClient),
		repository.NewHexRepository(entClient),
		repository.NewHexInfluenceRepository(entClient),
		repository.NewHexLeaderboardRepository(entClient),
		 logger),
		HexService: NewHexService(repository.NewHexRepository(entClient), logger),
		HexLeaderboardService: NewHexLeaderboardService(repository.NewHexLeaderboardRepository(entClient),
		repository.NewHexInfluenceRepository(entClient),
		 logger),
		HexInfluenceService: NewHexInfluenceService(repository.NewHexInfluenceRepository(entClient), logger),
	}
}

package handler

import (
	"stride-wars-app/internal/service"

	"go.uber.org/zap"
)

type Handlers struct {
	AuthHandler           *AuthHandler
	UserHandler           *UserHandler
	ActivityHandler       *ActivityHandler
	HexLeaderboardHandler *HexLeaderboardHandler
}

func Provide(services *service.Services, logger *zap.Logger) *Handlers {
	return &Handlers{
		AuthHandler:           NewAuthHandler(services.AuthService, logger),
		UserHandler:           NewUserHandler(services.UserService, logger),
		ActivityHandler:       NewActivityHandler(services.ActivityService, logger),
		HexLeaderboardHandler: NewHexLeaderboardHandler(services.HexLeaderboardService, logger),
	}
}

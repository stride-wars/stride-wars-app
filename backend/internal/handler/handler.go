package handler

import (
	"stride-wars-app/internal/service"

	"go.uber.org/zap"
)

type Handlers struct {
	AuthHandler    *AuthHandler
	HexDataHandler *HexDataHandler
}

func Provide(services *service.Services, logger *zap.Logger) *Handlers {
	return &Handlers{AuthHandler: NewAuthHandler(services.AuthService, logger)}
}

package handler

import (
	"go.uber.org/zap"
	"stride-wars-app/internal/service"
)

type Handlers struct {
	AuthHandler *AuthHandler
}

func Provide(services *service.Services, logger *zap.Logger) *Handlers {
	return &Handlers{AuthHandler: NewAuthHandler(services.AuthService, logger)}
}

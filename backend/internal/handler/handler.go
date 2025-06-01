package handler

import (
	"stride-wars-app/internal/service"

	"go.uber.org/zap"
	"stride-wars-app/ent"
)

type Handlers struct {
	AuthHandler    *AuthHandler
	HexDataHandler *HexDataHandler
}

func Provide(services *service.Services, logger *zap.Logger, entClient *ent.Client) *Handlers {
	return &Handlers{
		AuthHandler: NewAuthHandler(services.AuthService, logger),
		HexDataHandler: NewHexDataHandler(logger, entClient),
	}
}

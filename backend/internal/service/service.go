package service

import (
	"go.uber.org/zap"
	"stride-wars-app/ent"
	"stride-wars-app/internal/repository"
)

type Service struct {
	UserService *UserService
}

func Provide(client *ent.Client, logger *zap.Logger) *Service {
	return &Service{
		UserService: NewUserService(repository.NewUserRepository(client), logger),
	}
}

package service

import (
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
	"stride-wars-app/internal/repository"
)

type Services struct {
	UserService *UserService
	AuthService *AuthService
}

func Provide(repositories *repository.Repositories, supabaseClient *supabase.Client, logger *zap.Logger) *Services {
	userService := NewUserService(repositories.UserRepository, logger)

	return &Services{
		UserService: userService,
		AuthService: NewAuthService(supabaseClient, logger, userService),
	}
}

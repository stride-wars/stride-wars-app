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
}

func Provide(entClient *ent.Client, supabaseClient *supabase.Client, logger *zap.Logger) *Service {
	return &Service{
		UserService: NewUserService(repository.NewUserRepository(entClient), logger),
		AuthService: NewAuthService(supabaseClient, logger),
	}
}

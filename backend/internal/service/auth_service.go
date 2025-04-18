package service

import (
	"context"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

type AuthService struct {
	supabaseClient *supabase.Client
	logger         *zap.Logger
}

func NewAuthService(supabaseClient *supabase.Client, logger *zap.Logger) *AuthService {
	return &AuthService{
		supabaseClient: supabaseClient,
		logger:         logger,
	}
}

func (a *AuthService) SignUp(ctx context.Context, req types.SignupRequest) (interface{}, error) {
	_, err := a.supabaseClient.Auth.Signup(req)
	if err != nil {
		return nil, err
	}

	panic("implement me")
}

func (a *AuthService) SignIn(ctx context.Context, email string, password string) (*types.TokenResponse, error) {
	_, err := a.supabaseClient.Auth.SignInWithEmailPassword(email, password)
	if err != nil {
		return nil, err
	}

	panic("implement me")
}

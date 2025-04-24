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
	resp, err := a.supabaseClient.Auth.Signup(req)
	if err != nil {
		a.logger.Error("Failed to sign up", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func (a *AuthService) SignIn(ctx context.Context, email string, password string) (*types.TokenResponse, error) {
	resp, err := a.supabaseClient.Auth.SignInWithEmailPassword(email, password)
	if err != nil {
		a.logger.Error("Failed to sign in", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

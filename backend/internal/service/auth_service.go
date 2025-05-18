package service

import (
	"context"
	"stride-wars-app/ent/model"
	"stride-wars-app/pkg/errors"
	"stride-wars-app/pkg/utils"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

type AuthService struct {
	supabaseClient *supabase.Client
	logger         *zap.Logger
	userService    *UserService
}

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Session      types.Session `json:"session"`
	UserID       string        `json:"user_id"`
	ExternalUser string        `json:"external_user"`
	Username     string        `json:"username"`
	Email        string        `json:"email"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Session      types.Session `json:"session"`
	UserID       string        `json:"user_id"`
	ExternalUser string        `json:"external_user"`
	Username     string        `json:"username"`
	Email        string        `json:"email"`
}

func NewAuthService(supabaseClient *supabase.Client, logger *zap.Logger, userService *UserService) *AuthService {
	return &AuthService{
		supabaseClient: supabaseClient,
		logger:         logger,
		userService:    userService,
	}
}

func (a *AuthService) SignUp(ctx context.Context, req SignUpRequest) (*SignUpResponse, error) {
	err := a.validateSignUp(req)
	if err != nil {
		return nil, err
	}

	u, err := a.userService.FindByUsername(ctx, req.Username)
	if u != nil {
		return nil, errors.New("User already exists")
	}
	if err != nil {
		return nil, err
	}

	supabaseSignUp := types.SignupRequest{
		Email:         req.Email,
		Password:      "",
		Data:          nil,
		SecurityEmbed: types.SecurityEmbed{},
	}

	supabaseResp, err := a.supabaseClient.Auth.Signup(supabaseSignUp)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		ExternalUser: supabaseResp.User.ID,
	}
	internalUser, err := a.userService.CreateUser(ctx, user)
	if err != nil {
		a.logger.Error("Failed to sign up", zap.Error(err))
		return nil, err
	}

	resp := &SignUpResponse{
		Session:      supabaseResp.Session,
		UserID:       internalUser.ID.String(),
		Username:     internalUser.Username,
		ExternalUser: supabaseResp.User.ID.String(),
		Email:        supabaseResp.User.Email,
	}

	return resp, nil
}

func (a *AuthService) SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error) {
	supabaseResp, err := a.supabaseClient.Auth.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	internalUser, err := a.userService.FindByExternalUserID(ctx, supabaseResp.User.ID)
	if err != nil {
		a.logger.Error("Failed to sign in", zap.Error(err))
		return nil, err
	}

	resp := &SignInResponse{
		Session:  supabaseResp.Session,
		UserID:   internalUser.ID.String(),
		Username: internalUser.Username,
		Email:    supabaseResp.User.Email,
	}

	return resp, nil
}

func (a *AuthService) validateSignUp(req SignUpRequest) error {
	if !utils.IsValidEmail(req.Email) {
		return errors.New("Invalid email provided.")
	}

	if req.Username == "" {
		return errors.New("Username is required.")
	}

	if req.Password == "" {
		return errors.New("Password is required.")
	}

	return nil
}

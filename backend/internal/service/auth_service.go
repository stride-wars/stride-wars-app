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

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

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
	Data         string        `json:"data"`
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

	// Check if user exists
	u, err := a.userService.FindByUsername(ctx, req.Username)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if u != nil {
		return nil, errors.New("Username already exists!")
	}

	supabaseSignUp := types.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	supabaseResp, err := a.supabaseClient.Auth.Signup(supabaseSignUp)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		ExternalUser: supabaseResp.ID,
	}
	internalUser, err := a.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Create a new session by signing in after signup
	signInResp, err := a.supabaseClient.Auth.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	resp := &SignUpResponse{
		Session:      signInResp.Session,
		UserID:       internalUser.ID.String(),
		Username:     internalUser.Username,
		ExternalUser: supabaseResp.ID.String(),
		Email:        supabaseResp.Email,
	}

	return resp, nil
}

func (a *AuthService) SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error) {
	supabaseResp, err := a.supabaseClient.Auth.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		// Check if the error is due to unconfirmed email
		if err.Error() == "response status code 400: {\"code\":400,\"error_code\":\"email_not_confirmed\",\"msg\":\"Email not confirmed\"}" {
			return nil, errors.New("Please check your email for a confirmation link. If you haven't received it, try signing up again.")
		}

		return nil, err
	}

	supabaseUser := supabaseResp.User
	internalUser, err := a.userService.FindByExternalUserID(ctx, supabaseUser.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.New("User not found. Please sign up first.")
		}
		return nil, err
	}

	resp := &SignInResponse{
		Session:      supabaseResp.Session,
		UserID:       internalUser.ID.String(),
		Username:     internalUser.Username,
		ExternalUser: supabaseUser.ID.String(),
		Email:        supabaseUser.Email,
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

func (a *AuthService) ValidateToken(ctx context.Context, token string) (*Claims, error) {
	user, err := a.supabaseClient.Auth.GetUser()
	if err != nil {
		return nil, err
	}

	supabaseUser := user.User
	internalUser, err := a.userService.FindByExternalUserID(ctx, supabaseUser.ID)
	if err != nil {
		return nil, err
	}

	return &Claims{
		UserID: internalUser.ID.String(),
		Email:  supabaseUser.Email,
	}, nil
}

func (a *AuthService) ValidateSession(token string) (*types.User, error) {
	user, err := a.supabaseClient.Auth.GetUser()
	if err != nil {
		return nil, err
	}
	return &user.User, nil
}

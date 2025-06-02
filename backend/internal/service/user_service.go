package service

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"

	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UpdateUsernameRequest struct {
	OldUsername string `json:"old_username"`
	NewUsername string `json:"new_username"`
}

type UpdateUsernameResponse struct {
	ID          uuid.UUID `json:"id"`
	NewUsername string    `json:"new_username"`
}

type UserService struct {
	repository repository.UserRepository
	logger     *zap.Logger
}

func NewUserService(repository repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repository: repository, logger: logger}
}

func (us *UserService) FindByID(ctx context.Context, uuid uuid.UUID) (*ent.User, error) {
	return us.repository.FindByID(ctx, uuid)
}

func (us *UserService) FindByExternalUserID(ctx context.Context, uuid uuid.UUID) (*ent.User, error) {
	return us.repository.FindByExternalUserID(ctx, uuid)
}

func (us *UserService) FindByIDs(ctx context.Context, uuids []uuid.UUID) ([]*ent.User, error) {
	return us.repository.FindByIDs(ctx, uuids)
}

func (us *UserService) FindByUsername(ctx context.Context, username string) (*ent.User, error) {
	return us.repository.FindByUsername(ctx, username)
}

func (us *UserService) CreateUser(ctx context.Context, user *model.User) (*ent.User, error) {
	return us.repository.CreateUser(ctx, user)
}

// TO DO: update username

func (s *UserService) UpdateUsername(ctx context.Context, req *UpdateUsernameRequest) (*ent.User, error) {
	// 1) Fetch the existing user
	usr, err := s.repository.FindByUsername(ctx, req.OldUsername)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, err
	}

	// 2) Apply the change
	usr.Username = req.NewUsername

	// 3) Map ent.User to model.User and persist
	updatedUser := &model.User{
		ID:       usr.ID,
		Username: usr.Username,
		// Add other fields if necessary
	}
	rowsAffected, err := s.repository.UpdateUsername(ctx, updatedUser)
	if err != nil {
		return nil, err
	}
	if rowsAffected != 1 {
		return nil, err
	}

	// 4) Build response
	updatedUsr, err := s.repository.FindByUsername(ctx, req.NewUsername)
	if err != nil {
		return nil, err
	}
	return updatedUsr, nil
}

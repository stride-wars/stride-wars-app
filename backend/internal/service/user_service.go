package service

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"
)

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

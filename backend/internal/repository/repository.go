package repository

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"

	"github.com/google/uuid"
)

type Repositories struct {
	UserRepository     UserRepository
	ActivityRepository ActivityRepository
	//HexRepository *HexRepository
	//HexInfluenceRepository *HexInfluenceRepository
	//HexLeaderboardRepository *HexLeaderboardRepository
	//FriendshipRepository *FriendshipRepository
}

func Provide(client *ent.Client) *Repositories {
	return &Repositories{
		UserRepository:     NewUserRepository(client),
		ActivityRepository: NewActivityRepository(client),
	}
}

type IUserRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*ent.User, error)
	FindByIDs(ctx context.Context, uuids []uuid.UUID) ([]*ent.User, error)
	FindByExternalUserID(ctx context.Context, uuid uuid.UUID) (*ent.User, error)
	FindByUsername(ctx context.Context, username string) (*ent.User, error)
	CreateUser(ctx context.Context, user *model.User) (*ent.User, error)
	UpdateUsername(ctx context.Context, user *model.User) (int, error)
}

type IActivityRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*ent.Activity, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Activity, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.Activity, error)
	CreateActivity(ctx context.Context, activity *model.Activity) (*ent.Activity, error)
}

type IHexRepository interface{}

type IHexInfluence interface{}

type IHexLeaderboard interface{}

type IFriendship interface{}

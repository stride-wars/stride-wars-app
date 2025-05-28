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
	HexRepository HexRepository
	HexInfluenceRepository HexInfluenceRepository
	HexLeaderboardRepository HexLeaderboardRepository
	// FriendshipRepository *FriendshipRepository
}

func Provide(client *ent.Client) *Repositories {
	return &Repositories{
		UserRepository:     NewUserRepository(client),
		ActivityRepository: NewActivityRepository(client),
		HexRepository: NewHexRepository(client),
		HexInfluenceRepository: NewHexInfluenceRepository(client),
		HexLeaderboardRepository: NewHexLeaderboardRepository(client),
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

type IHexRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*ent.User, error)
	FindByIDs(ctx context.Context, uuids []uuid.UUID) ([]*ent.User, error)
	CreateHex(ctx context.Context, hex *model.Hex) (*ent.Hex, error)
	CreateHexes(ctx context.Context, hexes []*model.Hex) ([]*ent.Hex, error)
}

type IHexInfluenceRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*ent.Activity, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Activity, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.Activity, error)
	FindByHexID(ctx context.Context, h3_index int64) ([]*ent.Activity, error)
	FindByUserIDAndHexID(ctx context.Context, userID uuid.UUID, hexID int64) (*ent.HexInfluence, error)
	CreateHexInfluence(ctx context.Context, influence *model.HexInfluence) (*ent.HexInfluence, error)
	UpdateHexInfluence(ctx context.Context, influence *model.HexInfluence) (int, error)
	UpdateHexInfluences(ctx context.Context, userID uuid.UUID, hexIDs []int64) (int, error)
	UpdateOrCreateHexInfluence(ctx context.Context, userID uuid.UUID, hexID int64) (*ent.HexInfluence, error)
	UpdateOrCreateHexInfluences(ctx context.Context, userID uuid.UUID, hexIDs []int64) ([]*ent.HexInfluence, error)
}

type IHexLeaderboardRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*ent.HexLeaderboard, error)
	FindByH3Index(ctx context.Context, h3Index int64) (*ent.HexLeaderboard, error)
	FindByH3Indexes(ctx context.Context, h3Indexes []int64) ([]*ent.HexLeaderboard, error)
	CreateHexLeaderboard(ctx context.Context, hex *model.HexLeaderboard) (*ent.HexLeaderboard, error)
	UpdateHexLeaderboard(ctx context.Context, hex *model.HexLeaderboard) (int, error)
}

type IFriendship interface{}

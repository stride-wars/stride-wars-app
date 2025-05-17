package repository

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"

	"github.com/google/uuid"
)

type IUserRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*ent.User, error)
	FindByIDs(ctx context.Context, uuids []uuid.UUID) ([]*ent.User, error)
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
}

type IHexInfluenceRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*ent.Activity, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Activity, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.Activity, error)
	CreateHexInfluence(ctx context.Context, influence *model.HexInfluence) (*ent.HexInfluence, error)
	UpdateHexInfluence(ctx context.Context, influence *model.HexInfluence) (int, error)
}

type IHexLeaderboardRepository interface {
	FindByH3Index(ctx context.Context, h3Index int64) (*ent.HexLeaderboard, error)
	FindByH3Indexes(ctx context.Context, h3Indexes []int64) ([]*ent.HexLeaderboard, error)
	CreateHexLeaderboard(ctx context.Context, hex *model.HexLeaderboard) (*ent.HexLeaderboard, error)
	UpdateHexLeaderboard(ctx context.Context, hex *model.HexLeaderboard) (int, error)
}

type IFriendship interface{}

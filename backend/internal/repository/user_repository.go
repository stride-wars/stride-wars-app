package repository

import (
	"context"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	entUser "stride-wars-app/ent/user"

	"github.com/google/uuid"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return UserRepository{client: client}
}

func (r UserRepository) FindByID(ctx context.Context, uuid uuid.UUID) (*ent.User, error) {
	return r.client.User.Query().Where(entUser.IDEQ(uuid)).First(ctx)
}

func (r UserRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.User, error) {
	return r.client.User.Query().Where(entUser.IDIn(ids...)).All(ctx)
}

func (r UserRepository) FindByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().Where(entUser.UsernameEQ(username)).First(ctx)
}

func (r UserRepository) CreateUser(ctx context.Context, user *model.User) (*ent.User, error) {
	return r.client.User.Create().SetUsername(user.Username).SetExternalUser(user.ExternalUser).SetID(uuid.New()).Save(ctx)
}

func (r UserRepository) UpdateUsername(ctx context.Context, user *model.User) (int, error) {
	return r.client.User.Update().Where(entUser.IDEQ(user.ID)).SetUsername(user.Username).Save(ctx)
}

package repository

import (
	"context"
	"stride-wars-app/ent"
	entActivity "stride-wars-app/ent/activity"
	"stride-wars-app/ent/model"

	"github.com/google/uuid"
)

type ActivityRepository struct {
	client *ent.Client
}

func NewActivityRepository(client *ent.Client) IActivityRepository {
	return ActivityRepository{client: client}
}

func (r ActivityRepository) FindByID(ctx context.Context, uuid uuid.UUID) (*ent.Activity, error) {
	return r.client.Activity.Query().Where(entActivity.IDEQ(uuid)).First(ctx)
}

func (r ActivityRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Activity, error) {
	return r.client.Activity.Query().Where(entActivity.IDIn(ids...)).All(ctx)
}

func (r ActivityRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.Activity, error) {
	return r.client.Activity.Query().Where(entActivity.UserIDIn(userID)).All(ctx)
}

func (r ActivityRepository) CreateActivity(ctx context.Context, activity *model.Activity) (*ent.Activity, error) {
	return r.client.Activity.Create().SetID(uuid.New()).SetUserID(activity.UserID).SetDurationSeconds(activity.Duration).SetDistanceMeters(activity.Distance).SetH3Indexes(activity.H3Indexes).Save(ctx)
}

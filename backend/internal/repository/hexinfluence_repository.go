package repository

import (
	"context"
	"stride-wars-app/ent"
	entHexInfluence "stride-wars-app/ent/hexinfluence"
	"stride-wars-app/ent/model"

	"time"

	"github.com/google/uuid"
)

type HexInfluenceRepository struct 	{
	client *ent.Client
}
func NewHexInfluenceRepository(client *ent.Client) HexInfluenceRepository {
	return HexInfluenceRepository{client: client}
}
func (r HexInfluenceRepository) FindByID(ctx context.Context, uuid uuid.UUID) (*ent.HexInfluence, error) {
	return r.client.HexInfluence.Query().Where(entHexInfluence.IDEQ(uuid)).First(ctx)
}

func (r HexInfluenceRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.HexInfluence, error) {
	return r.client.HexInfluence.Query().Where(entHexInfluence.UserIDEQ(userID)).All(ctx)
}

func (r HexInfluenceRepository) FindByUserIDAndHexID(ctx context.Context, userID uuid.UUID, hexID int64) (*ent.HexInfluence, error) {
	return r.client.HexInfluence.Query().Where(
		entHexInfluence.H3IndexEQ(hexID),
		entHexInfluence.UserIDEQ(userID),
	).First(ctx)
}
func (r HexInfluenceRepository) FindByHexID(ctx context.Context, hexID int64) ([]*ent.HexInfluence, error) {
	return r.client.HexInfluence.Query().Where(entHexInfluence.H3IndexEQ(hexID)).All(ctx)
}

func (r HexInfluenceRepository) CreateHexInfluence(ctx context.Context, hexInfluence *model.HexInfluence) (*ent.HexInfluence, error) {
	return r.client.HexInfluence.Create().
		SetH3Index(hexInfluence.H3Index).
		SetUserID(hexInfluence.UserID).
		SetScore(hexInfluence.Score).
		SetLastUpdated(hexInfluence.LastUpdated).
		SetID(uuid.New()).
		Save(ctx)
}
func (r HexInfluenceRepository) UpdateHexInfluence(ctx context.Context,  userID uuid.UUID, hexID int64) (int, error) {
	hexInfluence, err := r.FindByUserIDAndHexID(ctx, userID, hexID)
	if err != nil {
		return 0, err
	}
	if hexInfluence == nil {
		return 0, nil
	}

	new_update_time := time.Now()
	multiplyer := 1 - 0.1*new_update_time.Sub(hexInfluence.LastUpdated).Hours() / 24.0 / 7.0
	multiplyer = float64(int(multiplyer*10)) / 10.0
	if multiplyer < 0.0 {
		multiplyer = 0.1
	} 

	new_score := hexInfluence.Score * multiplyer + 1.0

	return r.client.HexInfluence.Update().
		Where(entHexInfluence.IDEQ(hexInfluence.ID)).
		SetLastUpdated(new_update_time).
		SetScore(new_score).
		Save(ctx)
}

func (r HexInfluenceRepository) UpdateHexInfluences(ctx context.Context,userID uuid.UUID,hexIDs []int64,) (int, error) {
    totalUpdated := 0
    for _, h3id := range hexIDs {
        n, err := r.UpdateHexInfluence(ctx, userID, h3id)
        if err != nil {
            return totalUpdated, err
        }
            totalUpdated += n
        }
    
    return totalUpdated, nil
}


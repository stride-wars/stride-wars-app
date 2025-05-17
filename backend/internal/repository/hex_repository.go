package repository

import (
	"context"
	"stride-wars-app/ent"
	entHex "stride-wars-app/ent/hex"
)

type HexRepository struct {
	client *ent.Client
}

func NewHexRepository(client *ent.Client) HexRepository {
	return HexRepository{client: client}
}

func (r HexRepository) FindByID(ctx context.Context, hex_id int64) (*ent.Hex, error) {
	return r.client.Hex.Query().Where(entHex.IDEQ(hex_id)).First(ctx)
}

func (r HexRepository) FindByIDs(ctx context.Context, ids []int64) ([]*ent.Hex, error) {
	return r.client.Hex.Query().Where(entHex.IDIn(ids...)).All(ctx)
}

func (r HexRepository) CreateHex(ctx context.Context, h3_index int64) (*ent.Hex, error) {
	return r.client.Hex.Create().SetID(h3_index).Save(ctx)
}

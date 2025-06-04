package model

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type HexInfluence struct {
	ID          uuid.UUID
	H3Index     int64
	UserID      uuid.UUID
	Score       float64
	LastUpdated time.Time
	ent.Schema
}

func (HexInfluence) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.Int64("h3_index"),
		field.UUID("user_id", uuid.UUID{}),
		field.Float("score"),
		field.Time("last_updated"),
	}
}

func (HexInfluence) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("hex", Hex.Type).Field("h3_index").Unique().Required(),
		edge.To("users", User.Type).Unique().Field("user_id").Required(),
	}
}

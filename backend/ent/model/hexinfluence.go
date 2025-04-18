package model

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type HexInfluence struct {
	ent.Schema
}

func (HexInfluence) Fields() []ent.Field {
	return []ent.Field{
		field.String("h3_index"),
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

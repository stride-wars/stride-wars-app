package model

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type HexLeaderboard struct {
	ent.Schema
}

func (HexLeaderboard) Fields() []ent.Field {
	return []ent.Field{
		field.String("h3_index").Unique(),
		field.JSON("top_users", map[string][]uuid.UUID{}),
	}
}

func (HexLeaderboard) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("hex", Hex.Type).Field("h3_index").Unique().Required(),
	}
}

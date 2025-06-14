package model

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type TopUser struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Score    float64   `json:"score"`
}

type HexLeaderboard struct {
	ID       uuid.UUID
	H3Index  string
	TopUsers []TopUser `json:"top_users"`
	ent.Schema
}

func (HexLeaderboard) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("h3_index").Unique(),
		field.JSON("top_users", []TopUser{}).
			Default([]TopUser{}),
	}
}

func (HexLeaderboard) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("hex", Hex.Type).Field("h3_index").Unique().Required(),
	}
}

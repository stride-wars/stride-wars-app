package model

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Friendship struct {
	ent.Schema
}

func (Friendship) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("friend_id", uuid.UUID{}),
		field.Time("created_at"),
	}
}

func (Friendship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type).Unique().Field("user_id").Required(),
		edge.To("friends", User.Type).Unique().Field("friend_id").Required(),
	}
}

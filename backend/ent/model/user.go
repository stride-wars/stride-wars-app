package model

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"

	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	ExternalUser uuid.UUID
	Username     string
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("external_user", uuid.UUID{}),
		field.String("username"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("activity", Activity.Type).Ref("users"),
		edge.From("friendship", Friendship.Type).Ref("users"),
		edge.From("hexinfluence", HexInfluence.Type).Ref("users"),
	}
}

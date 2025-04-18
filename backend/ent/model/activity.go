package model

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

type Activity struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Duration  float64
	Distance  float64
	H3Indexes []string
	CreatedAt time.Time
	ent.Schema
}

func (Activity) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.Float("duration_seconds"),
		field.Float("distance_meters"),
		field.Strings("h3_indexes"),
		field.Time("created_at"),
	}
}

func (Activity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type).Unique().Field("user_id").Required(),
	}
}

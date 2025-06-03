package model

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
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
		field.JSON("h3_indexes", []string{}),
		field.Time("created_at").Default(time.Now()),
	}
}

func (Activity) Edges() []ent.Edge {
	return []ent.Edge{
		// This *creates* the FK on activity.user_id → user.id
		edge.To("user", User.Type).
			Field("user_id").
			Unique().
			Required(),
	}
}

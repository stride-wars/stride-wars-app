package model

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Hex struct {
	ent.Schema
}

func (Hex) ID() ent.Field {
	return field.Int64("id").
		Immutable()
}

func (Hex) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Immutable(),
	}
}

func (Hex) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("hexinfluences", HexInfluence.Type).Ref("hex"),
		edge.From("hexleaderboards", HexLeaderboard.Type).Ref("hex"),
	}
}

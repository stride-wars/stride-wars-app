package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Hex struct {
	ent.Schema
}

// TODO implement proper fields and edges when schema is ready
func (Hex) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("hex_owner"),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.Bool("is_active"),
	}
}

func (Hex) Edges() []ent.Edge {
	return []ent.Edge{}
}

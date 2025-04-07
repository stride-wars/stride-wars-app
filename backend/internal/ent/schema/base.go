package schema

import "entgo.io/ent"

// Base holds the schema definition for the Base entity.
type Base struct {
	ent.Schema
}

// Fields of the Base.
func (Base) Fields() []ent.Field {
	return nil
}

// Edges of the Base.
func (Base) Edges() []ent.Edge {
	return nil
}

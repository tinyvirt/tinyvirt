package schema

import "entgo.io/ent"

// VM holds the schema definition for the VM entity.
type VM struct {
	ent.Schema
}

// Fields of the VM.
func (VM) Fields() []ent.Field {
	return nil
}

// Edges of the VM.
func (VM) Edges() []ent.Edge {
	return nil
}

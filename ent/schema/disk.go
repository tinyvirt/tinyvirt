package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Disk holds the schema definition for the Disk entity.
type Disk struct {
	ent.Schema
}

// Fields of the Disk.
func (Disk) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("name").NotEmpty().Unique(),
		field.String("format").NotEmpty(),
		field.String("description").Optional(),
		field.Uint32("size_gb").Validate(func(value uint32) error {
			if value < 1 {
				return fmt.Errorf("disk size must be at least 1 GB")
			}
			return nil
		}),
	}
}

// Edges of the Disk.
func (Disk) Edges() []ent.Edge {
	return nil
}

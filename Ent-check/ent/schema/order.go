// ent/schema/order.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		// Nếu bạn không dùng ID mặc định (int auto-increment), hãy disable nó trong `ID()` method.
		field.Uint64("order_id").
			Unique().
			Comment("Primary key - auto increment"),

		field.Int("version").
			Comment("Version for concurrency"),

		field.String("note").
			MaxLen(200),

		field.Time("created_at"),

		field.Time("updated_at").
			Nillable().
			Optional(),

		field.String("created_by").
			MaxLen(100),

		field.String("updated_by").
			MaxLen(100).
			Nillable().
			Optional(),
	}
}

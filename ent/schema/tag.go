package schema

import (
	"cp-website/util"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Tag struct {
	ent.Schema
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(util.GenID).Immutable(),
		field.String("name").NotEmpty().Unique(),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		// 反向关联到 CP
		edge.From("cps", CP.Type).Ref("tags"),
	}
}

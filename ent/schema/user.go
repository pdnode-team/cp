package schema

import (
	"cp-website/util"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(util.GenID).Immutable(),
		// 保存 Zitadel 的 sub，必须唯一
		field.String("sub").Unique().NotEmpty(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// 一个用户可以拥有多个 CP 和 Tag
		edge.To("cps", CP.Type),
		edge.To("tags", Tag.Type),
		edge.To("liked_cps", CP.Type),
		edge.To("comments", Comment.Type),
	}
}

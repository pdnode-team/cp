package schema

import (
	"cp-website/util"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type CP struct {
	ent.Schema
}

func (CP) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(util.GenID).Immutable(),
		// 对应你之前的校验：至少1个字符，唯一
		field.String("name").NotEmpty().Unique(),
		// 对应之前的 2-20 个字符
		field.String("category").MinLen(2).MaxLen(20),
		// 对应可选的 link
		field.String("link").Optional(),
	}
}

func (CP) Edges() []ent.Edge {
	return []ent.Edge{
		// 定义到 Tag 的多对多关系，Ent 会自动生成中间表
		edge.To("tags", Tag.Type),
		edge.From("owner", User.Type).Ref("cps").Unique().Required(),
	}
}
